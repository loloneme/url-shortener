package persistence

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type UrlRepository interface {
	Save(ctx context.Context, model *domain.ShortenedURL) error
	GetByOriginal(ctx context.Context, original string) (*domain.ShortenedURL, error)
	GetByShort(ctx context.Context, short string) (*domain.ShortenedURL, error)
}
