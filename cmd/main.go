package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo"
	"github.com/nandanugg/prometheus-grafana-go-webserver-docker-example/config"
	"github.com/nandanugg/prometheus-grafana-go-webserver-docker-example/entity"
	"github.com/nandanugg/prometheus-grafana-go-webserver-docker-example/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestProm = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

func main() {
	cfg, err := config.LoadEnv()
	checkError(err)

	dbConn, err := createDBConnection(cfg.DBConfig)
	checkError(err)

	service := service.NewUserService(dbConn)

	// Echo example
	e := echo.New()
	e.Use(EchoPrometheusMiddleware)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	e.GET("/user", func(c echo.Context) error {
		err := hasError()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		users, err := service.GetAllUser(context.Background())
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, users)
	})

	type CreateUserRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	e.POST("/user", func(c echo.Context) error {
		err := hasError()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		req := new(CreateUserRequest)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = service.PostUser(context.Background(), entity.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "User created")
	})
	type UpdateUserRequest struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	e.PATCH("/user", func(c echo.Context) error {
		err := hasError()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		req := new(UpdateUserRequest)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = service.UpdateUserById(context.Background(), req.ID, entity.User{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "User updated")
	})

	type DeleteUserRequest struct {
		ID int `json:"id"`
	}

	e.DELETE("/user", func(c echo.Context) error {
		err := hasError()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		req := new(DeleteUserRequest)
		if err := c.Bind(req); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		err = service.DeleteUserById(context.Background(), req.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, "User deleted")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func createDBConnection(cfg config.DBConfig) (*pgx.Conn, error) {
	connURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Params,
	)

	conn, err := pgx.Connect(context.Background(), connURL)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func hasError() error {
	if rand.Float64() < 0.1 {
		return errors.New("an error occurred")
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Echo example
func EchoPrometheusMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c) // Process request

		status := c.Response().Status
		httpRequestProm.WithLabelValues(c.Request().URL.Path, c.Request().Method, fmt.Sprintf("%v", status)).Observe(float64(time.Since(start).Milliseconds()))

		return err
	}
}
