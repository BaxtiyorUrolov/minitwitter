package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
	"os"
)

type Config struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	HTTPort string

	EmailPassword string
	KafkaHost     string
	KafkaPort     string
	KafkaTopic    string

	ServiceName string
	LoggerLevel string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("error !", err.Error())
	}
	cfg := Config{}

	cfg.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "localhost"))
	cfg.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", "5432"))
	cfg.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "your user"))
	cfg.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "your password"))
	cfg.PostgresDB = cast.ToString(getOrReturnDefault("POSTGRES_DB", "your database"))
	cfg.HTTPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	cfg.EmailPassword = cast.ToString(getOrReturnDefault("EMAIL_PASSWORD", "your password"))
	cfg.KafkaHost = cast.ToString(getOrReturnDefault("KAFKA_HOST", "localhost"))
	cfg.KafkaPort = cast.ToString(getOrReturnDefault("KAFKA_PORT", "9092"))
	cfg.KafkaTopic = cast.ToString(getOrReturnDefault("KAFKA_TOPIC", "topic"))
	cfg.ServiceName = cast.ToString(getOrReturnDefault("SERVICE_NAME", "store"))
	cfg.LoggerLevel = cast.ToString(getOrReturnDefault("LOGGER_LEVEL", "debug"))
	return cfg
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	return defaultValue
}
