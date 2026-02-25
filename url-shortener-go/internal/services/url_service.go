package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"url-shortener-go/internal/models"
	"url-shortener-go/internal/repositories"
	"url-shortener-go/pkg/random"

	"github.com/redis/go-redis/v9"
)

type URLService struct {
	repo  repositories.URLRepository
	redis *redis.Client
}

func NewURLService(repo repositories.URLRepository, redisURL models.RedisURL) *URLService {
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     redisCfg.Address,
	// 	Password: redisCfg.Password,
	// 	DB:       redisCfg.DB,
	// })
	opts, err := redis.ParseURL(redisURL.URL)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse REDIS_URL: %v", err))
	}

	rdb := redis.NewClient(opts)

	return &URLService{
		repo:  repo,
		redis: rdb,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, originalURL string, expiry time.Duration) (string, error) {
	// Генерация короткого кода
	shortCode, err := random.GenerateShortCode(originalURL)
	if err != nil {
		return "", err
	}

	// Создание новой записи URL в базе данных.
	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(expiry),
	}

	// Сохранение в базе данных
	if err := s.repo.Create(ctx, url); err != nil {
		return "", fmt.Errorf("fail to save URL: %w", err)
	}

	// Кэширование в Redis
	if err := s.cacheURL(ctx, url); err != nil {
		fmt.Printf("Warning: failed to cache URL: %v\n", err)
	}

	return shortCode, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	// Попытка получить из кэша
	cachedURL, err := s.redis.Get(ctx, shortCode).Result()
	if err == nil {
		return cachedURL, err
	}

	// Пытаемся получить из базы данных
	url, err := s.repo.FindByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("Error: Url not found : %w", err)
	}

	// Кэширование результата для будущих запросов
	if err := s.cacheURL(ctx, url); err != nil {
		fmt.Printf("Warning: failed to cache URL: %v\n", err)
	}

	return url.OriginalURL, nil
}

func (s *URLService) cacheURL(ctx context.Context, url *models.URL) error {
	ttl := time.Until(url.ExpiresAt)

	if ttl <= 0 {
		return errors.New("URL already expried")
	}

	return s.redis.Set(ctx, url.ShortCode, url.OriginalURL, ttl).Err()
}
