package shortenedurls

import (
	"context"
	"sync"
	domain "url-shortener/internal/domain/shortenedurl"
)

type Model = domain.ShortenedURL

type Repository struct {
	mu sync.RWMutex

	byShort    map[string]string
	byOriginal map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		byShort:    make(map[string]string),
		byOriginal: make(map[string]string),
	}
}

func (r *Repository) Save(ctx context.Context, model *Model) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byOriginal[model.Original]; exists {
		return nil
	}

	if _, exists := r.byShort[model.Short]; exists {
		return domain.ErrDuplicate
	}

	r.byOriginal[model.Original] = model.Short
	r.byShort[model.Short] = model.Original

	return nil
}

func (r *Repository) GetByOriginal(ctx context.Context, original string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	short, exists := r.byOriginal[original]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return &Model{
		Short:    short,
		Original: original,
	}, nil
}

func (r *Repository) GetByShort(ctx context.Context, short string) (*Model, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, exists := r.byShort[short]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return &Model{
		Short:    short,
		Original: original,
	}, nil
}
