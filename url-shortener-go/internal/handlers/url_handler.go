package handlers

import (
	"net/http"
	"time"

	"url-shortener-go/internal/models"
	"url-shortener-go/internal/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

// POST
func (h *URLHandler) CreateURL(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Дефолтное время жизни  30 дней
	if req.ExpiresIn == nil || *req.ExpiresIn <= 0 {
		defaultExpiry := 720 * time.Hour // 30 days
		req.ExpiresIn = &defaultExpiry
	}

	shortCode, err := h.service.ShortenURL(c.Request.Context(), req.URL, *req.ExpiresIn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Образование полного короткого URL для ответа клиенту
	shortURL := "http://" + c.Request.Host + "/" + shortCode
	c.JSON(http.StatusCreated, models.ShortenResponse{ShortURL: shortURL})
}

// GET
func (h *URLHandler) GetURL(c *gin.Context) {
	shortCode := c.Param("shortCode")

	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}
