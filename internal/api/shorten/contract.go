//go:generate go run go.uber.org/mock/mockgen@latest -destination=./mocks/contract.go -package=mock -source=contract.go

package shorten

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type ShortenUrlService interface {
	ShortenUrl(ctx context.Context, url string) (*domain.ShortenedURL, error)
}

type shortenUrlRequest struct {
	Url string `json:"url"`
}

type shortenUrlResponse struct {
	ShortUrl string `json:"short_url"`
}
