package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/auth/domain"
)

type authRouterCtx struct {
	apiKeyUsecase domain.APIKeyUsecase
}

func NewAuthRouter(e *echo.Echo, apiKeyUsecase domain.APIKeyUsecase) {
	routerGroup := e.Group("/auth")
	routerCtx := &authRouterCtx{
		apiKeyUsecase,
	}

	routerGroup.POST("/keys", routerCtx.CreateAPIKey)
	routerGroup.DELETE("/keys/:publicID", routerCtx.RevokeAPIKey)
}

func (ctx *authRouterCtx) CreateAPIKey(c echo.Context) error {
	var apikey APIKeyIn

	if err := c.Bind(&apikey); err != nil {
		return HTTPError{Message: "Invalid request"}.ErrUnprocessableEntity() // TODO: Change to BadRequest()
	}

	// TODO: Validate the request body

	apiKey, err := ctx.apiKeyUsecase.CreateAPIKey(c.Request().Context(), apikey.Name)
	if err != nil {
		return HTTPError{Message: "Failed to create API key"}.InternalServerError()
	}

	apikeyOut := &APIKeyOut{
		ID:        apiKey.PublicID,
		Name:      apiKey.Name,
		Key:       apiKey.Key,
		CreatedAt: apiKey.CreatedAt,
	}

	return c.JSON(http.StatusCreated, apikeyOut)
}

func (r *authRouterCtx) RevokeAPIKey(c echo.Context) error {
	return nil
}
