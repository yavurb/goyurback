package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/projects/domain"
)

type projectRouterCtx struct {
	projectUsecase domain.ProjectUsecase
}

func NewProjectsRouter(e *echo.Echo, projectUsecase domain.ProjectUsecase) *projectRouterCtx {
	routerGroup := e.Group("/projects")
	routerCtx := &projectRouterCtx{
		projectUsecase,
	}

	routerGroup.POST("", routerCtx.createProject)
	routerGroup.GET("", routerCtx.getProjects)
	routerGroup.GET("/:id", routerCtx.getProject)

	return routerCtx
}

func (ctx *projectRouterCtx) createProject(c echo.Context) error {
	var project ProjectIn

	if err := c.Bind(&project); err != nil {
		c.Logger().Error(err)

		return HTTPError{
			Message: "Invalid request body",
		}.ErrUnprocessableEntity()
	}

	project_, err := ctx.projectUsecase.Create(
		c.Request().Context(), project.Name, project.Description, project.ThumbnailURL, project.WebsiteURL, project.Live, project.Tags, project.PostId,
	)
	if err != nil {
		return handleErr(err)
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

func (ctx *projectRouterCtx) getProject(c echo.Context) error {
	var params GetProjectParam

	if err := c.Bind(&params); err != nil {
		return HTTPError{
			Message: "Invalid params",
		}.BadRequest()
	}

	project_, err := ctx.projectUsecase.Get(c.Request().Context(), params.ID)
	if err != nil {
		return handleErr(err)
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

	return c.JSON(http.StatusOK, projectOut)
}

func (ctx *projectRouterCtx) getProjects(c echo.Context) error {
	projects, err := ctx.projectUsecase.GetProjects(c.Request().Context())
	if err != nil {
		return handleErr(err)
	}

	projectsOut := []*ProjectOut{}

	for _, project := range projects {
		projectsOut = append(projectsOut, &ProjectOut{
			ID:           project.PublicID,
			Name:         project.Name,
			Description:  project.Description,
			Tags:         project.Tags,
			ThumbnailURL: project.ThumbnailURL,
			WebsiteURL:   project.WebsiteURL,
			Live:         project.Live,
			CreatedAt:    project.CreatedAt,
			UpdatedAt:    project.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, &ProjectsOut{
		Data: projectsOut,
	})
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
