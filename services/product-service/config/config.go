package config

import (
	"fmt"
	"net/url"
	"os"
)

type AppConfiguration struct {
	AppPort string
	DB      DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func (c DBConfig) URL() string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.User, c.Password),
		Host:   c.Host + ":" + c.Port,
		Path:   c.Name,
	}
	q := url.Values{}
	u.RawQuery = q.Encode()
	return u.String()
}

func Load() (*AppConfiguration, error) {
	//Read secrets from environment variables
	cfg := &AppConfiguration{
		AppPort: getEnv("APP_PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			Name:     getEnv("DB_NAME", "product_db"),
			User:     getEnv("DB_USER", "product_user"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}

	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be set")
	}
	return cfg, nil
}

func getEnv(envVar string, defaultVal string) string {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		return v
	}
	return defaultVal
}
