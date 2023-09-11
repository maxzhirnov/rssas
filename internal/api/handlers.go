package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type handlers struct {
}

func newHandlers() *handlers {
	return &handlers{}
}

func (h handlers) ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
