package config

import (
	"github.com/Netflix/go-env"
)

// DBConfig represents the configuration structure for database variables.
type DBConfig struct {
	Name     string `env:"DB_NAME,required=true"`
	Port     string `env:"DB_PORT,required=true"`
	Host     string `env:"DB_HOST,required=true"`
	Username string `env:"DB_USERNAME,required=true"`
	Password string `env:"DB_PASSWORD,required=true"`
	Params   string `env:"DB_PARAMS,required=true"`
}

// EnvConfig represents the configuration structure for environment variables.
type EnvConfig struct {
	DBConfig DBConfig `env:""`
	Env      string   `env:"ENV,default=development"`
}

func LoadEnv() (*EnvConfig, error) {
	// Load environment variables using Netflix/go-env
	var cfg EnvConfig
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		// Handle error, e.g., log it or return a default configuration
		return nil, err
	}

	return &cfg, nil
}
