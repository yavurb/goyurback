package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/posts/infrastructure/ui"
)

type appContext struct {
	Settings *appSetings
}
type appSetings struct {
	Port string
}

func NewAppContext() *appContext {
	ctx := &appContext{
		Settings: &appSetings{},
	}

	ctx.initAppSettings()

	return ctx
}

func (c *appContext) NewRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true

	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "Healthy!") })

	ui.NewPostsRouter(e)

	return e
}

func (c *appContext) initAppSettings() {
	goenv := "dev"

	if value, ok := os.LookupEnv("GO_ENV"); ok {
		goenv = value
	}

	envsFileName := fmt.Sprintf(".env.%s", goenv)

	envs, err := godotenv.Read(envsFileName)
	if err != nil {
		log.Fatalf("Error loading `%s` file. Make sure the file is present and it is free of errors", envsFileName)
	}

	c.Settings.Port = envs["PORT"]
}
