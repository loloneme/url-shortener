//go:generate go run go.uber.org/mock/mockgen@latest -destination=./mocks/contract.go -package=mock -source=contract.go

package shorten

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type repo interface {
	Save(ctx context.Context, url *domain.ShortenedURL) error
	GetByOriginal(ctx context.Context, original string) (*domain.ShortenedURL, error)
}

type generator interface {
	Generate() string
}
