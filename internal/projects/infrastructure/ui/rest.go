package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

type projectRouterCtx struct {
	projectUsecase domain.ProjectUsecase
}

func NewProjectsRouter(e *echo.Echo, projectUsecase domain.ProjectUsecase) {
	routerGroup := e.Group("/projects")
	routerCtx := &projectRouterCtx{
		projectUsecase,
	}

	routerGroup.POST("", routerCtx.createProject)
}

func (ctx *projectRouterCtx) createProject(c echo.Context) error {
	var project ProjectIn

	if err := c.Bind(&project); err != nil {
		return HTTPError{
			Message: "Invalid request body",
		}.ErrUnprocessableEntity()
	}

	project_, err := ctx.projectUsecase.Create(
		c.Request().Context(), project.Name, project.Description, project.ThumbnailURL, project.WebsiteURL, project.Live, project.Tags, project.PostId,
	)

	if err != nil {
		handleErr(err)
	}

	projectOut := &ProjectOut{
		ID:           project_.PublicID,
		Name:         project_.Name,
		Description:  project_.Description,
		Tags:         project_.Tags,
		ThumbnailURL: project_.ThumbnailURL,
		WebsiteURL:   project_.WebsiteURL,
		Live:         project_.Live,
		CreatedAt:    project_.CreatedAt,
		UpdatedAt:    project_.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, projectOut)
}

func handleErr(err error) error {
	switch err {
	case domain.ErrProjectNotFound:
		return HTTPError{
			Message: "Project not found",
		}.NotFound()
	default:
		return HTTPError{
			Message: "Internal server error",
		}.InternalServerError()
	}
}
