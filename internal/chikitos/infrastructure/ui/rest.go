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

func NewChikitosRouter(e *echo.Echo, usecase domain.ChikitoUsecase) *chikitoRouterCtx {
	routerGroup := e.Group("/chikitos")
	routerCtx := &chikitoRouterCtx{
		usecase: usecase,
	}

	routerGroup.POST("", routerCtx.create)
	routerGroup.GET("/:id", routerCtx.get)

	return routerCtx
}

func (ctx *chikitoRouterCtx) create(c echo.Context) error {
	chikito := new(CreateIn)

	if err := c.Bind(chikito); err != nil {
		log.Printf("Bad request body for a chikito. %v", err)

		return HTTPError{
			Message: "Bad request body",
		}.ErrUnprocessableEntity()
	}

	if err := c.Validate(chikito); err != nil {
		// TODO: Format field errors and return a more helpful response message
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

func (ctx *chikitoRouterCtx) get(c echo.Context) error {
	var chikitoParams GetChikitoParams

	// TODO: Use fluent binding
	if err := c.Bind(&chikitoParams); err != nil {
		log.Printf("Bad chikito params. %v\n", err)

		return HTTPError{
			Message: "Bad chikito params",
		}.ErrUnprocessableEntity()
	}

	if err := c.Validate(chikitoParams); err != nil {
		// TODO: Format field errors and return a more helpful response message
		return HTTPError{
			Message: "Bad request params",
		}.ErrUnprocessableEntity()
	}

	chikito, err := ctx.usecase.Get(c.Request().Context(), chikitoParams.ID)
	if err != nil {
		log.Printf("Could not get chikito. %v\n", err)

		return HTTPError{
			Message: "Unable to get chikito",
		}.NotFound()
	}

	return c.Redirect(http.StatusPermanentRedirect, chikito.URL)
}
