package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	AppEnv         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	JWTSecret      string
	JWTExpiryHours int
	RabbitMQURL    string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	hours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	App = &Config{
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		AppEnv:         getEnv("APP_ENV", "development"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "rcpuser"),
		DBPassword:     getEnv("DB_PASSWORD", "rcppassword"),
		DBName:         getEnv("DB_NAME", "rcpdb"),
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpiryHours: hours,
		RabbitMQURL:    getEnv("RABBITMQ_URL", "amqp://rcpuser:rcppassword@localhost:5672/"),
	}

	log.Println("Config loaded successfully")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}