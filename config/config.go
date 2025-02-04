package config

import (
	"os"
)

const (
	dbConnectionString = "DB_CONNECTION_STRING"
	jwtSecret          = "JWT_SECRET"
)

// Config contains application configuration
type Config struct {
	DBConnectionString string
	JWTSecret          string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}
	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}
	// default configuration
	config := &Config{
		DBConnectionString: getEnvOrDefault(dbConnectionString, "postgres://user:password@localhost/account_db?sslmode=disable"),
		JWTSecret:          getEnvOrDefault(jwtSecret, "verysecrettext"),
	}

	return config, nil
}
