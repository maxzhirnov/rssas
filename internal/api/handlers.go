package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"rssas/internal/service"
)

type handlers struct {
	app *service.App
}

func newHandlers(app *service.App) *handlers {
	return &handlers{
		app: app,
	}
}

func (h handlers) ping(c echo.Context) error {
	responseData := map[string]string{
		"message": "Pong",
	}
	return c.JSON(http.StatusOK, responseData)
}

type AddFeedRequest struct {
	FeedURL string `json:"feed_url"`
}

func (h handlers) addFeed(c echo.Context) error {
	addFeed := new(AddFeedRequest)
	if err := c.Bind(addFeed); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.app.AddNewFeed(addFeed.FeedURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, addFeed)
}
