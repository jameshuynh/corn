package config

import "github.com/labstack/echo"

// ApplicationController is a controller that order controller should inherit
type ApplicationController interface {
	BeforeAction(string, string, echo.Context) error
	RenderJSON(interface{}) error
}
