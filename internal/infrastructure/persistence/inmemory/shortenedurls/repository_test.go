package shortenedurls

import (
	"context"
	"fmt"
	"sync"
	"testing"
	domain "url-shortener/internal/domain/shortenedurl"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepositorySave(t *testing.T) {
	t.Run("saves and returns model with generated short", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		model, err := repo.Save(ctx, "https://example.com")

		require.NoError(t, err)
		assert.Equal(t, uint64(1), model.ID)
		assert.Equal(t, "https://example.com", model.Original)
		assert.Len(t, model.Short, 10)
	})

	t.Run("returns existing model for duplicate original", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		first, err := repo.Save(ctx, "https://example.com")
		require.NoError(t, err)

		second, err := repo.Save(ctx, "https://example.com")
		require.NoError(t, err)

		assert.Equal(t, first.ID, second.ID)
		assert.Equal(t, first.Short, second.Short)
	})

	t.Run("increments ID for different originals", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		first, _ := repo.Save(ctx, "https://example-1.com")
		second, _ := repo.Save(ctx, "https://example-2.com")

		assert.Equal(t, uint64(1), first.ID)
		assert.Equal(t, uint64(2), second.ID)
		assert.NotEqual(t, first.Short, second.Short)
	})
}

func TestRepositoryGetByShort(t *testing.T) {
	t.Run("returns saved URL by short", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		saved, err := repo.Save(ctx, "https://example.com")
		require.NoError(t, err)

		result, err := repo.GetByShort(ctx, saved.Short)

		require.NoError(t, err)
		assert.Equal(t, "https://example.com", result.Original)
		assert.Equal(t, saved.ID, result.ID)
	})

	t.Run("returns not found for unknown short", func(t *testing.T) {
		repo := NewRepository()

		_, err := repo.GetByShort(context.Background(), "abcdefghij")

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestRepositoryGetByOriginal(t *testing.T) {
	t.Run("returns saved URL by original", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		saved, err := repo.Save(ctx, "https://example.com")
		require.NoError(t, err)

		result, err := repo.GetByOriginal(ctx, "https://example.com")

		require.NoError(t, err)
		assert.Equal(t, saved.Short, result.Short)
		assert.Equal(t, saved.ID, result.ID)
	})

	t.Run("returns not found for unknown original", func(t *testing.T) {
		repo := NewRepository()

		_, err := repo.GetByOriginal(context.Background(), "https://example.com")

		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestConcurrentAccess(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := range goroutines {
		go func(n int) {
			defer wg.Done()
			original := fmt.Sprintf("https://example.com/%d", n)
			_, _ = repo.Save(ctx, original)
		}(i)
	}
	wg.Wait()

	for i := range goroutines {
		original := fmt.Sprintf("https://example.com/%d", i)
		result, err := repo.GetByOriginal(ctx, original)
		require.NoError(t, err)
		assert.Equal(t, original, result.Original)

		byShort, err := repo.GetByShort(ctx, result.Short)
		require.NoError(t, err)
		assert.Equal(t, original, byShort.Original)
	}
}
