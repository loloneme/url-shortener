package shortenedurls

import (
	"context"
	"fmt"
	"sync"
	domain "url-shortener/internal/domain/shortenedurl"
	"url-shortener/internal/domain/shortgen"
)

type Model = domain.ShortenedURL

type Repository struct {
	mu         sync.RWMutex
	counter    uint64
	byShort    map[string]Model
	byOriginal map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		byShort:    make(map[string]Model),
		byOriginal: make(map[string]string),
	}
}

func (r *Repository) Save(_ context.Context, original string) (*Model, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if short, exists := r.byOriginal[original]; exists {
		model := r.byShort[short]
		return &model, nil
	}

	r.counter++
	short, err := shortgen.EncodeID(r.counter)
	if err != nil {
		return nil, fmt.Errorf("failed to encode id: %w", err)
	}

	model := Model{ID: r.counter, Short: short, Original: original}
	r.byShort[short] = model
	r.byOriginal[original] = short

	return &model, nil
}

func (r *Repository) GetByOriginal(_ context.Context, original string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	short, exists := r.byOriginal[original]
	if !exists {
		return nil, domain.ErrNotFound
	}

	model := r.byShort[short]
	return &model, nil
}

func (r *Repository) GetByShort(_ context.Context, short string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	model, exists := r.byShort[short]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return &model, nil
}
