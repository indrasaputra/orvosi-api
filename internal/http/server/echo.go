package server

import (
	"github.com/indrasaputra/orvosi-api/internal/http/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server acts as echo.Echo server.
type Server struct {
	*echo.Echo
}

// NewServer creates an instance of Echo.
func NewServer(jwtDecoder echo.MiddlewareFunc, routes []*router.Route) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	for _, route := range routes {
		var midds []echo.MiddlewareFunc
		midds = append(midds, jwtDecoder)
		midds = append(midds, route.Middlewares...)
		e.Add(route.Method, route.Path, route.Handler, midds...)
	}

	return &Server{e}
}
