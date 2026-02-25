package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"url-shortener-go/internal/models"
)

type URLRepository interface {
	Create(ctx context.Context, url *models.URL) error
	FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error)
}

type postgresURLRepository struct {
	db *pgxpool.Pool
}

// Создание новой записи URL в базе данных.
func (p *postgresURLRepository) Create(ctx context.Context, url *models.URL) error {
	query := `INSERT INTO urls (short_code, original_url, created_at, expires_at)
				VALUES ($1, $2, $3, $4)
	`

	_, err := p.db.Exec(ctx, query,
		url.ShortCode,
		url.OriginalURL,
		url.CreatedAt,
		url.ExpiresAt)

	return err
}

// FindByShortCode реализует URLRepository и извлекает запись URL из базы данных по короткому коду.
func (p *postgresURLRepository) FindByShortCode(ctx context.Context, shortCode string) (*models.URL, error) {
	query := `SELECT short_code, original_url, created_at, expires_at
			FROM urls
			WHERE short_code=$1
	`
	var url models.URL
	err := p.db.QueryRow(ctx, query, shortCode).Scan(
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.ExpiresAt)

	return &url, err
}

func NewURLRepository(cfg models.PostgresURL) (URLRepository, error) {
	pool, err := pgxpool.New(context.Background(), cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &postgresURLRepository{db: pool}, nil
}
