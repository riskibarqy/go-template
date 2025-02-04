package config

import (
	"log"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

const (
	appMode            = "APP_MODE"
	dbConnectionString = "DB_CONNECTION_STRING"
	jwtSecret          = "JWT_SECRET"
	redisAddr          = "REDIS_ADDR"
	redisPassword      = "REDIS_PASSWORD"
)

// Config contains application configuration
type Config struct {
	AppMode            string `json:"appMode"`
	DBConnectionString string `json:"dbConnectionString"`
	JWTSecret          string `json:"jwtSecret"`
	RedisAddr          string `json:"redisAddr"`
	RedisPassword      string `json:"redisPassword"`

	DatabaseClient *sqlx.DB
}

type Metadata struct {
	RedisExpirationShort  int `json:"redisExpirationShort"`
	RedisExpirationMedium int `json:"redisExpirationMedium"`
	RedisExpirationLong   int `json:"redisExpirationLong"`
}

var AppConfig = &Config{}
var MetadataConfig = &Metadata{}

// getEnvOrDefault retrieves the value of an environment variable or returns a default value.
// It supports string, int, bool, and float64 types.
func getEnvOrDefault(env string, defaultVal interface{}) interface{} {
	e, _ := os.LookupEnv(env)
	if e == "" {
		return defaultVal
	}

	switch v := defaultVal.(type) {
	case string:
		return e
	case int:
		if intVal, err := strconv.Atoi(e); err == nil {
			return intVal
		}
		return v // return default if conversion fails
	case bool:
		if boolVal, err := strconv.ParseBool(e); err == nil {
			return boolVal
		}
		return v // return default if conversion fails
	case float64:
		if floatVal, err := strconv.ParseFloat(e, 64); err == nil {
			return floatVal
		}
		return v // return default if conversion fails
	default:
		return defaultVal // return default for unsupported types
	}
}

// GetConfiguration retrieves application configuration based on set environment
func GetConfiguration() {
	// Load .env file in development
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	AppConfig.AppMode = getEnvOrDefault(appMode, "development").(string)
	AppConfig.DBConnectionString = getEnvOrDefault(dbConnectionString, "postgres://user:password@localhost/account_db?sslmode=disable").(string)
	AppConfig.JWTSecret = getEnvOrDefault(jwtSecret, "verysecrettext").(string)
	AppConfig.RedisAddr = getEnvOrDefault(redisAddr, "localhost:6379").(string)
	AppConfig.RedisPassword = getEnvOrDefault(redisPassword, "password").(string)
}
