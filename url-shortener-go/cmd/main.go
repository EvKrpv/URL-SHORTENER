package main

import (
	"log"

	"url-shortener-go/internal/config"
	"url-shortener-go/internal/handlers"
	"url-shortener-go/internal/repositories"
	"url-shortener-go/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации из .env файла и переменных окружения
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Инициализация репозитория
	repo, err := repositories.NewURLRepository(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Инициализация сервисов
	urlService := services.NewURLService(repo, cfg.Redis)

	// Инициализация обработчиков
	urlHandler := handlers.NewURLHandler(urlService)

	// Инициализация маршрутизатора
	r := gin.Default()

	// маршруты
	r.POST("/shorten", urlHandler.CreateURL)
	r.GET("/:shortCode", urlHandler.GetURL)

	// запуск сервера
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
