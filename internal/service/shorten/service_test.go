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
	t.Run("retries on duplicate and succeeds", testRetriesOnDuplicate)
}

func TestShortenServiceFail(t *testing.T) {
	t.Run("repo error on GetByOriginal", testRepoGetError)
	t.Run("repo error on Save", testRepoSaveError)
	t.Run("fails after max retries", testFailsAfterMaxRetries)
}

func testReturnsExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	existing := &domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}
	mockRepo.EXPECT().
		GetByOriginal(ctx, existing.Original).
		Return(existing, nil)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, existing.Original)

	require.NoError(t, err)
	assert.Equal(t, existing, result)
}

func testSavesNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	mockGen.EXPECT().
		Generate().
		Return("abcdefghij")

	mockRepo.EXPECT().
		Save(ctx, &domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}).
		Return(nil)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, "https://example.com")

	require.NoError(t, err)
	assert.Equal(t, "abcdefghij", result.Short)
	assert.Equal(t, "https://example.com", result.Original)
}

func testRetriesOnDuplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	gomock.InOrder(
		mockGen.EXPECT().Generate().Return("collision1"),
		mockGen.EXPECT().Generate().Return("collision2"),
		mockGen.EXPECT().Generate().Return("unique999"),
	)

	gomock.InOrder(
		mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(domain.ErrDuplicate),
		mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(domain.ErrDuplicate),
		mockRepo.EXPECT().Save(ctx, gomock.Any()).Return(nil),
	)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, "https://example.com")

	require.NoError(t, err)
	assert.Equal(t, "unique999", result.Short)
}

func testFailsAfterMaxRetries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	mockGen.EXPECT().
		Generate().
		Return("collision").
		Times(maxRetries)

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(domain.ErrDuplicate).
		Times(maxRetries)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, "https://example.com")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, domain.ErrDuplicate)
	assert.Contains(t, err.Error(), "failed to generate unique short after 5 attempts")
}

func testRepoGetError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	repoErr := errors.New("database connection error")
	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, repoErr)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, "https://example.com")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

func testRepoSaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	mockGen := mock.NewMockgenerator(ctrl)
	ctx := context.Background()

	saveErr := errors.New("database write error")
	mockRepo.EXPECT().
		GetByOriginal(ctx, "https://example.com").
		Return(nil, domain.ErrNotFound)

	mockGen.EXPECT().
		Generate().
		Return("abcdefghij")

	mockRepo.EXPECT().
		Save(ctx, gomock.Any()).
		Return(saveErr)

	service := New(mockRepo, mockGen)
	result, err := service.ShortenUrl(ctx, "https://example.com")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, saveErr)
}
