package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"rssas/internal/service"
)

type Server struct {
	app      *service.App
	echo     *echo.Echo
	handlers *handlers
}

func NewServer(app *service.App) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return &Server{
		app:      app,
		echo:     e,
		handlers: newHandlers(),
	}
}

func (s Server) Run(address string) error {
	s.echo.GET("/ping", s.handlers.ping)
	if err := s.echo.Start(address); err != nil {
		return err
	}
	log.Infof("Starting server on %s", address)
	return nil
}
