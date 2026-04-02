package shorten

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	mock "url-shortener/internal/api/shorten/mocks"
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

func TestShortenHandlerSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockShortenUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.EXPECT().
		ShortenUrl(gomock.Any(), "https://example.com").
		Return(&domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}, true, nil)

	handler := New(mockService)
	err := handler.ShortenUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.JSONEq(t, `{"short_url":"http://localhost:8080/abcdefghij/redirect"}`, rec.Body.String())
}

func TestShortenHandlerReturnsExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockShortenUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.EXPECT().
		ShortenUrl(gomock.Any(), "https://example.com").
		Return(&domain.ShortenedURL{Short: "abcdefghij", Original: "https://example.com"}, false, nil)

	handler := New(mockService)
	err := handler.ShortenUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"short_url":"http://localhost:8080/abcdefghij/redirect"}`, rec.Body.String())
}

func TestShortenUrlHandlerFail(t *testing.T) {
	t.Run("bad request on invalid request body", testBadRequestInvalidBody)
	t.Run("bad request on empty url", testBadRequestEmptyUrl)
	t.Run("bad request on url without protocol", testBadRequestNoProtocol)
	t.Run("internal server error", testInternalServerError)
}

func testBadRequestInvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockShortenUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`not json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := New(mockService)
	err := handler.ShortenUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func testBadRequestEmptyUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockShortenUrlService(ctrl)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := New(mockService)
	err := handler.ShortenUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func testBadRequestNoProtocol(t *testing.T) {
	invalidUrls := []string{
		"example.com",
		"ftp://example.com",
		"just-text",
	}

	for _, url := range invalidUrls {
		t.Run(url, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockShortenUrlService(ctrl)

			e := echo.New()
			body := `{"url":"` + url + `"}`
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := New(mockService)
			err := handler.ShortenUrl(c)

			require.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		})
	}
}

func testInternalServerError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockShortenUrlService(ctrl)
	serviceError := errors.New("service error")

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(`{"url":"https://example.com"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService.EXPECT().
		ShortenUrl(gomock.Any(), "https://example.com").
		Return(nil, false, serviceError)

	handler := New(mockService)
	err := handler.ShortenUrl(c)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}
