package shorten

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"url-shortener/internal/infrastructure/logger"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ShortenUrlService
}

func New(service ShortenUrlService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ShortenUrl(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	httpReq := new(shortenUrlRequest)
	if err := c.Bind(httpReq); err != nil {
		log.Warn("ShortenUrl invalid request", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.validateRequest(httpReq); err != nil {
		log.Warn("ShortenUrl invalid request", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	urlModel, created, err := h.service.ShortenUrl(c.Request().Context(), httpReq.Url)
	if err != nil {
		log.Error("ShortenUrl failed", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}

	return c.JSON(status, shortenUrlResponse{ShortUrl: fmt.Sprintf("http://localhost:8080/%s/redirect", urlModel.Short)})
}

func (h *Handler) validateRequest(shortenUrlRequest *shortenUrlRequest) error {
	if shortenUrlRequest.Url == "" {
		return errors.New("url is required")
	}

	if !strings.HasPrefix(shortenUrlRequest.Url, "http://") && !strings.HasPrefix(shortenUrlRequest.Url, "https://") {
		return errors.New("url must start with http:// or https://")
	}

	return nil
}
