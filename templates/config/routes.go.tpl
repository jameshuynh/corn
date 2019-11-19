package config

import (
	"net/http"

	"github.com/labstack/echo"
)

// Routes setups all the routing
func Routes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusCreated, "API Server is online")
	})
}
