package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	postApplication "github.com/yavurb/goyurback/internal/posts/application"
	postRepository "github.com/yavurb/goyurback/internal/posts/infrastructure/repository"
	postUI "github.com/yavurb/goyurback/internal/posts/infrastructure/ui"

	projectApplication "github.com/yavurb/goyurback/internal/projects/application"
	projectRepository "github.com/yavurb/goyurback/internal/projects/infrastructure/repository"
	projectUI "github.com/yavurb/goyurback/internal/projects/infrastructure/ui"

	authApplication "github.com/yavurb/goyurback/internal/auth/application"
	authRepository "github.com/yavurb/goyurback/internal/auth/infrastructure/repository"
	authUI "github.com/yavurb/goyurback/internal/auth/infrastructure/ui"
)

type appContext struct {
	Settings *appSetings
	Connpool *pgxpool.Pool
	ctx      context.Context
}
type appSetings struct {
	Port         string
	DBConnString string
}

func NewAppContext() *appContext {
	appCtx := &appContext{
		Settings: &appSetings{},
		ctx:      context.Background(),
	}

	appCtx.initAppSettings()

	connpool, err := pgxpool.New(appCtx.ctx, appCtx.Settings.DBConnString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	appCtx.Connpool = connpool

	return appCtx
}

// TODO: Add custom validator
func (c *appContext) NewRouter() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // Limits request to 20req/s based on the client's IP
	e.Use(middleware.Logger())                                              // Use a simple logger middleware

	e.GET("/health", func(c echo.Context) error { return c.String(http.StatusOK, "Healthy!") })

	postRespository := postRepository.NewRepo(c.Connpool)
	postUcase := postApplication.NewPostUsecase(postRespository)
	postUI.NewPostsRouter(e, postUcase)

	projectRespository := projectRepository.NewRepo(c.Connpool)
	projectUcase := projectApplication.NewProjectUsecase(projectRespository)
	projectUI.NewProjectsRouter(e, projectUcase)

	authAPIKeyRespository := authRepository.NewAPIKeyRepo(c.Connpool)
	authAPIKeyUcase := authApplication.NewAPIKeyUsecase(authAPIKeyRespository)
	authUI.NewAuthRouter(e, authAPIKeyUcase)

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-api-key",
		Validator: func(auth string, c echo.Context) (bool, error) {
			isValid, err := authAPIKeyUcase.ValidateAPIKey(c.Request().Context(), auth)

			return isValid, err
		},
	}))

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
	c.Settings.DBConnString = envs["DB_URI"]
}
