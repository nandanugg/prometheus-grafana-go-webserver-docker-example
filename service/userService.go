package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nandanugg/prometheus-grafana-go-webserver-docker-example/entity"
)

type UserService struct {
	db *pgx.Conn
}

func NewUserService(db *pgx.Conn) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetAllUser(ctx context.Context) ([]entity.User, error) {
	rows, err := s.db.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *UserService) PostUser(ctx context.Context, user entity.User) error {
	_, err := s.db.Exec(ctx, "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserById(ctx context.Context, id int, user entity.User) error {
	_, err := s.db.Exec(ctx, "UPDATE users SET username = $1, password = $2 WHERE id = $3", user.Username, user.Password, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUserById(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
