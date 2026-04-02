package redirect

import (
	"errors"
	"net/http"
	domain "url-shortener/internal/domain/shortenedurl"
	"url-shortener/internal/domain/shortgen"
	"url-shortener/internal/infrastructure/logger"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service getOriginalUrlService
}

func New(service getOriginalUrlService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RedirectToOriginal(c echo.Context) error {
	log := logger.FromContext(c.Request().Context())

	short := c.Param("short")
	if err := shortgen.Validate(short); err != nil {
		log.Warn("RedirectToOriginal invalid request", "error", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	urlModel, err := h.service.GetOriginalUrl(c.Request().Context(), short)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("RedirectToOriginal not found", "error", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "URL not found"})
		}
		log.Error("RedirectToOriginal failed", "error", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	log.Info("RedirectToOriginal success", "original_url", urlModel.Original)
	return c.Redirect(http.StatusFound, urlModel.Original)
}
