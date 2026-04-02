//go:generate go run go.uber.org/mock/mockgen@latest -destination=./mocks/contract.go -package=mock -source=contract.go

package getoriginal

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type getOriginalUrlService interface {
	GetOriginalUrl(ctx context.Context, shortUrl string) (*domain.ShortenedURL, error)
}

type getOriginalUrlResponse struct {
	OriginalUrl string `json:"original_url"`
}
