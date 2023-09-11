package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"rssas/internal/log"
	"rssas/internal/service"
)

type handlers struct {
	app    *service.App
	logger *log.Logger
}

func newHandlers(app *service.App, logger *log.Logger) *handlers {
	return &handlers{
		app:    app,
		logger: logger,
	}
}

func (h handlers) ping(c echo.Context) error {
	h.logger.Log.Info("ping handler fired")
	responseData := map[string]string{
		"message": "Pong",
	}
	return c.JSON(http.StatusOK, responseData)
}

type AddFeedRequest struct {
	FeedURL string `json:"feed_url"`
}

func (h handlers) addFeed(c echo.Context) error {
	h.logger.Log.Info("addFeed handler fired")
	addFeed := new(AddFeedRequest)
	if err := c.Bind(addFeed); err != nil {
		h.logger.Log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := h.app.AddNewFeed(addFeed.FeedURL)
	if err != nil {
		h.logger.Log.Error(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, addFeed)
}
