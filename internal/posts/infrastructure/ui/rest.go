package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type postRouterCtx struct {
	e *echo.Echo
}

func NewPostsRouter(e *echo.Echo) {
	routerGroup := e.Group("/posts")
	routerCtx := &postRouterCtx{
		e: e,
	}

	routerGroup.POST("", routerCtx.getPosts)
}

func (ctx *postRouterCtx) getPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, &struct{ Hello string }{Hello: "world"})
}
