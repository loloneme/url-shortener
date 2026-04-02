package shorten

import (
	"context"
	"errors"
	"fmt"
	domain "url-shortener/internal/domain/shortenedurl"
)

type Service struct {
	repo repo
}

func New(repo repo) *Service {
	return &Service{repo: repo}
}

func (s *Service) ShortenUrl(ctx context.Context, original string) (model *domain.ShortenedURL, created bool, err error) {
	model, err = s.repo.GetByOriginal(ctx, original)
	if err == nil {
		return model, false, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, false, fmt.Errorf("failed to get by original: %w", err)
	}

	model, err = s.repo.Save(ctx, original)
	if err != nil {
		return nil, false, fmt.Errorf("failed to save shortened url: %w", err)
	}

	return model, true, nil
}
