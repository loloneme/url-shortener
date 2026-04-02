package shorten

import (
	"context"
	"errors"
	"fmt"
	domain "url-shortener/internal/domain/shortenedurl"
)

const maxRetries = 5

type Service struct {
	repo      repo
	generator generator
}

func New(repo repo, generator generator) *Service {
	return &Service{repo: repo, generator: generator}
}

func (s *Service) ShortenUrl(ctx context.Context, original string) (*domain.ShortenedURL, error) {
	model, err := s.repo.GetByOriginal(ctx, original)
	if err == nil {
		return model, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		return nil, fmt.Errorf("failed to get by original: %w", err)
	}

	var lastErr error
	for range maxRetries {
		generatedShort := s.generator.Generate()
		newURL := &domain.ShortenedURL{
			Short:    generatedShort,
			Original: original,
		}
		lastErr = s.repo.Save(ctx, newURL)
		if lastErr == nil {
			return newURL, nil
		}
		if !errors.Is(lastErr, domain.ErrDuplicate) {
			return nil, fmt.Errorf("failed to save shortened url: %w", lastErr)
		}
	}
	return nil, fmt.Errorf("failed to generate unique short after %d attempts: %w", maxRetries, lastErr)
}
