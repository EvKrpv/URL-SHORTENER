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
	// load env variagle
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Initialize repository
	repo, err := repositories.NewURLRepository(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize services
	urlService := services.NewURLService(repo, cfg.Redis)

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(urlService)

	// step a router
	r := gin.Default()

	// routes
	r.POST("/shorten", urlHandler.CreateURL)
	r.GET("/:shortCode", urlHandler.GetURL)

	// starting a server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
