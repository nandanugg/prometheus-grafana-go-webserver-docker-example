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

// AuthConfig represents the configuration structure for authentication variables.
type AuthConfig struct {
	JWTSecret  string `env:"JWT_SECRET,required=true"`
	BCryptSalt int    `env:"BCRYPT_SALT,required=true,default=8"`
}

// S3Config represents the configuration structure for S3 variables.
type S3Config struct {
	ID         string `env:"S3_ID,required=true"`
	SecretKey  string `env:"S3_SECRET_KEY,required=true"`
	BucketName string `env:"S3_BUCKET_NAME,required=true"`
	Region     string `env:"S3_REGION,required=true"`
}

// EnvConfig represents the configuration structure for environment variables.
type EnvConfig struct {
	DBConfig   DBConfig   `env:""`
	AuthConfig AuthConfig `env:""`
	S3Config   S3Config   `env:""`
	Env        string     `env:"ENV,default=development"`
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
