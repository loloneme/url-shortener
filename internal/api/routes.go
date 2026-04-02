package api

import (
	"url-shortener/internal/api/getoriginal"
	"url-shortener/internal/api/redirect"
	"url-shortener/internal/api/shorten"

	"github.com/labstack/echo/v4"
)

type API struct {
	shortenHandler     *shorten.Handler
	redirectHandler    *redirect.Handler
	getOriginalHandler *getoriginal.Handler
}

func NewAPI(
	shortenHandler *shorten.Handler,
	redirectHandler *redirect.Handler,
	getOriginalHandler *getoriginal.Handler,
) *API {
	return &API{
		shortenHandler:     shortenHandler,
		redirectHandler:    redirectHandler,
		getOriginalHandler: getOriginalHandler,
	}
}

func (a *API) InitRoutes(e *echo.Echo) {
	e.POST("/api/shorten", a.shortenHandler.ShortenUrl)
	e.GET("/:short/redirect", a.redirectHandler.RedirectToOriginal)
	e.GET("/api/:short", a.getOriginalHandler.GetOriginalUrl)
}
