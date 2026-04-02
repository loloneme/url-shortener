package shorten

import (
	"context"
	"errors"
	"testing"
	domain "url-shortener/internal/domain/shortenedurl"
	mock "url-shortener/internal/service/shorten/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestShortenServiceSuccess(t *testing.T) {
	t.Run("returns existing URL", testReturnsExisting)
	t.Run("saves new URL", testSavesNew)
}

func TestShortenServiceFail(t *testing.T) {
	t.Run("repo error on GetByOriginal", testRepoGetError)
	t.Run("repo error on Save", testRepoSaveError)
}

func testReturnsExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	existing := &domain.ShortenedURL{ID: 1, Short: "aaaaaaaaab", Original: "https://example.com"}
	mockRepo.EXPECT().
		GetByOriginal(ctx, existing.Original).
		Return(existing, nil)

	service := New(mockRepo)
	result, created, err := service.ShortenUrl(ctx, existing.Original)

	require.NoError(t, err)
	assert.False(t, created)
	assert.Equal(t, existing, result)
}

func testSavesNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	saved := &domain.ShortenedURL{ID: 1, Short: "aaaaaaaaab", Original: "https://example.com"}
	mockRepo.EXPECT().
		Save(ctx, "https://example.com").
		Return(saved, nil)

	service := New(mockRepo)
	result, created, err := service.ShortenUrl(ctx, "https://example.com")

	require.NoError(t, err)
	assert.True(t, created)
	assert.Equal(t, saved, result)
}

func testRepoGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	repoErr := errors.New("database connection error")
	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, repoErr)

	service := New(mockRepo)
	result, created, err := service.ShortenUrl(ctx, "https://example.com")

	require.Error(t, err)
	assert.False(t, created)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

func testRepoSaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	saveErr := errors.New("database write error")
	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	mockRepo.EXPECT().
		Save(ctx, "https://example.com").
		Return(nil, saveErr)

	service := New(mockRepo)
	result, created, err := service.ShortenUrl(ctx, "https://example.com")

	require.Error(t, err)
	assert.False(t, created)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, saveErr)
}
