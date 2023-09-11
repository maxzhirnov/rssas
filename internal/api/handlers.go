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
	responseData := map[string]string{
		"message": "Pong!",
	}
	return c.JSON(http.StatusOK, responseData)
}
