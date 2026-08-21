package config

import (
	"expense-tracker/internal/utils"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}
type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret                string
	JwtExpiresIn          time.Duration
	RefreshTokenExpiresIn time.Duration
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(utils.GetEnv("JWT_EXPIRES_IN", "15m"))
	refreshTokenExpiresIn, _ := time.ParseDuration(utils.GetEnv("REFRESH_TOKEN_EXPIRES_IN", "15m"))

	config := &Config{
		Server: ServerConfig{
			Port:    utils.GetEnv("PORT", "8080"),
			GinMode: utils.GetEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			Name:     utils.GetEnv("DB_NAME", "ecommerce"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:                utils.GetEnv("JWT_SECRET", "your_jwt_secret"),
			JwtExpiresIn:          jwtExpiresIn,
			RefreshTokenExpiresIn: refreshTokenExpiresIn,
		},
	}

	return config, nil
}
