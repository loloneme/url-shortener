//go:generate go run go.uber.org/mock/mockgen@latest -destination=./mocks/contract.go -package=mock -source=contract.go

package redirect

import (
	"context"
	domain "url-shortener/internal/domain/shortenedurl"
)

type getOriginalUrlService interface {
	GetOriginalUrl(ctx context.Context, short string) (*domain.ShortenedURL, error)
}
