package app

import "github.com/labstack/echo/v4"

type appContext struct{}

func NewAppContext() *appContext {
	return &appContext{}
}

func (c *appContext) NewRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true

	return e
}
