package getoriginal

import (
	"context"
	"errors"
	"testing"
	domain "url-shortener/internal/domain/shortenedurl"
	mock "url-shortener/internal/service/getoriginal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetOriginalSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	existing := &domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}

	mockRepo.EXPECT().
		GetByShort(ctx, existing.Short).
		Return(existing, nil)

	service := New(mockRepo)

	result, err := service.GetOriginalUrl(ctx, existing.Short)
	require.NoError(t, err)
	assert.Equal(t, existing, result)
}

func TestGetOriginalRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockrepo(ctrl)
	ctx := context.Background()

	existing := &domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}

	repoErr := errors.New("database connection error")
	mockRepo.EXPECT().
		GetByShort(ctx, existing.Short).
		Return(existing, repoErr)

	service := New(mockRepo)

	result, err := service.GetOriginalUrl(ctx, existing.Short)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}
