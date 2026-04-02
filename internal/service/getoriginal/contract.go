//go:generate go run go.uber.org/mock/mockgen@latest -destination=./mocks/contract.go -package=mock -source=contract.go

package getoriginal

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type repo interface {
	GetByShort(ctx context.Context, short string) (*domain.ShortenedURL, error)
}
