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
	t.Run("saves and returns existing URL by short", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		err := repo.Save(ctx, &Model{Short: "abcdefghij", Original: "https://example.com"})
		require.NoError(t, err)

		result, err := repo.GetByShort(ctx, "abcdefghij")

		require.NoError(t, err)
		assert.Equal(t, "https://example.com", result.Original)
	})

	t.Run("saves and returns existing URL by original", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()

		err := repo.Save(ctx, &Model{Short: "abcdefghij", Original: "https://example.com"})
		require.NoError(t, err)

		result, err := repo.GetByOriginal(ctx, "https://example.com")

		require.NoError(t, err)
		assert.Equal(t, "abcdefghij", result.Short)
	})

	t.Run("duplicate short returns ErrDuplicate", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()
		_ = repo.Save(ctx, &Model{Short: "abcdefghij", Original: "https://example-1.com"})
		err := repo.Save(ctx, &Model{Short: "abcdefghij", Original: "https://example-2.com"})

		assert.ErrorIs(t, err, domain.ErrDuplicate)
	})

	t.Run("duplicate original is idempotent", func(t *testing.T) {
		repo := NewRepository()
		ctx := context.Background()
		_ = repo.Save(ctx, &Model{Short: "abcdefghij", Original: "https://example.com"})
		err := repo.Save(ctx, &Model{Short: "1234567890", Original: "https://example.com"})
		assert.NoError(t, err)
		result, _ := repo.GetByOriginal(ctx, "https://example.com")
		assert.Equal(t, "abcdefghij", result.Short, "original short should not be overwritten")
	})
}

func TestRepositoryGetByShort(t *testing.T) {
	t.Run("return not found error if short is not found", func(t *testing.T) {
		repo := NewRepository()
		_, err := repo.GetByShort(context.Background(), "abcdefghij")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}

func TestRepositoryGetByOriginal(t *testing.T) {
	t.Run("return not found error if original is not found", func(t *testing.T) {
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
			short := fmt.Sprintf("short_%04d", n)
			original := fmt.Sprintf("https://example.com/%d", n)
			_ = repo.Save(ctx, &Model{Short: short, Original: original})
		}(i)
	}
	wg.Wait()
	for i := range goroutines {
		short := fmt.Sprintf("short_%04d", i)
		result, err := repo.GetByShort(ctx, short)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("https://example.com/%d", i), result.Original)
	}
}
