package config

import (
	"log"
	"os"
	"time"

	"url-shortener-go/internal/models"
	"github.com/joho/godotenv"
)

func Load() (*models.Config, error) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("error loading .env file")
	}

	return &models.Config{
		Postgres: models.PostgresURL{
			URL: getEnv("DB_URL", "postgres://user:securepassword@localhost:5432/table_name?sslmode=disable"),
		},
		Redis: models.RedisURL{
			URL: getEnv("REDIS_URL", "redis://localhost:6379/0"),
		},
		Server: models.ServerConfig{
			Port:        getEnv("PORT", "8080"),
			Environment: getEnv("ENVIRONMENT", "production"),
			Timeout:     30 * time.Second,
		},
	}, err
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Println("Retrive the env variables")
		return value
	}

	return defaultValue
}
