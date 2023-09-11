package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"rssas/internal/log"
	"rssas/internal/service"
)

type Server struct {
	app      *service.App
	echo     *echo.Echo
	handlers *handlers
	logger   *log.Logger
}

func NewServer(app *service.App, logger *log.Logger) *Server {
	e := echo.New()
	e.Use(middleware.Recover())
	return &Server{
		app:      app,
		echo:     e,
		handlers: newHandlers(app, logger),
		logger:   logger,
	}
}

func (s Server) Run(address string) error {
	s.echo.GET("/ping", s.handlers.ping)
	s.echo.POST("/add-feed", s.handlers.addFeed)
	if err := s.echo.Start(address); err != nil {
		return err
	}
	s.logger.Log.Infof("Starting server on %s", address)
	return nil
}
