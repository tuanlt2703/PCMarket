package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	HTTPPort string
	DB       DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
	MaxConns int32
}

func (c DBConfig) URL() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   c.Name,
	}
	q := url.Values{}
	q.Set("sslmode", c.SSLMode)
	q.Set("pool_max_conns", strconv.Itoa(int(c.MaxConns)))
	u.RawQuery = q.Encode()
	return u.String()
}

func Load() (*Config, error) {
	maxConns, err := strconv.Atoi(getEnv("DB_MAX_CONNS", "10"))
	if err != nil || maxConns <= 0 {
		return nil, fmt.Errorf("DB_MAX_CONNS must be a positive integer")
	}

	cfg := &Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "product_db"),
			User:     getEnv("DB_USER", "product_user"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: int32(maxConns),
		},
	}

	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be set")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
