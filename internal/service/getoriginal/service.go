package getoriginal

import (
	"context"
	"fmt"
	domain "url-shortener/internal/domain/shortenedurl"
)

type Service struct {
	repo repo
}

func New(repo repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetOriginalUrl(ctx context.Context, short string) (*domain.ShortenedURL, error) {
	urlModel, err := s.repo.GetByShort(ctx, short)
	if err != nil {
		return nil, fmt.Errorf("failed to get original url: %w", err)
	}
	return urlModel, nil
}
