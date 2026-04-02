package getoriginal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	mock "url-shortener/internal/api/getoriginal/mocks"
	domain "url-shortener/internal/domain/shortenedurl"
	"url-shortener/internal/infrastructure/logger"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func init() {
	logger.Init()
}

func TestGetOriginalHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockgetOriginalUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/abcdefghij", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short")
	c.SetParamValues("abcdefghij")

	mockService.EXPECT().
		GetOriginalUrl(gomock.Any(), "abcdefghij").
		Return(&domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}, nil)

	handler := New(mockService)
	err := handler.GetOriginalUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"original_url":"https://example.com"}`, rec.Body.String())
}

func TestGetOriginalUrlHandlerFail(t *testing.T) {
	t.Run("invalid request", testBadRequest)
	t.Run("original url not found", testNotFound)
	t.Run("internal error", testInternalServerError)
}

func testBadRequest(t *testing.T) {
	invalidShorts := []string{"", "abc", "abcde!ghij", "abcdefghijk", "abc/defghi"}

	for _, short := range invalidShorts {
		t.Run(short, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockgetOriginalUrlService(ctrl)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/api/"+short, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("short")
			c.SetParamValues(short)

			handler := New(mockService)
			err := handler.GetOriginalUrl(c)

			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func testNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockgetOriginalUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/abcdefghij", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short")
	c.SetParamValues("abcdefghij")

	mockService.EXPECT().
		GetOriginalUrl(gomock.Any(), "abcdefghij").
		Return(nil, domain.ErrNotFound)

	handler := New(mockService)
	err := handler.GetOriginalUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func testInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockgetOriginalUrlService(ctrl)
	serviceError := errors.New("service error")

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/abcdefghij", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("short")
	c.SetParamValues("abcdefghij")

	mockService.EXPECT().
		GetOriginalUrl(gomock.Any(), "abcdefghij").
		Return(nil, serviceError)

	handler := New(mockService)
	err := handler.GetOriginalUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
