package ui

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/chikitos/domain"
)

type chikitoRouterCtx struct {
	usecase domain.ChikitoUsecase
}

func NewChikitosRouter(echo *echo.Echo, usecase domain.ChikitoUsecase) {
	routerGroup := echo.Group("/chikito")
	routerCtx := &chikitoRouterCtx{
		usecase: usecase,
	}

	routerGroup.POST("", routerCtx.create)
}

func (ctx *chikitoRouterCtx) create(c echo.Context) error {
	var chikito CreateIn

	if err := c.Bind(&chikito); err != nil {
		log.Printf("Bad request body for a chikito. %v", err)

		return HTTPError{
			Message: "Bad request body",
		}.ErrUnprocessableEntity()
	}

	chikito_, err := ctx.usecase.Create(c.Request().Context(), chikito.URL, chikito.Description)
	if err != nil {
		log.Printf("Could not create chikito. %v", err)

		return HTTPError{
			Message: "Unable to create chikito",
		}.InternalServerError()
	}

	chikitoOut := &CreateOut{
		ID:          chikito_.PublicID,
		URL:         chikito_.URL,
		Description: chikito.Description,
		CreatedAt:   chikito_.CreatedAt,
		UpdatedAt:   chikito_.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, chikitoOut)
}
