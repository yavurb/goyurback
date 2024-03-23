package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yavurb/goyurback/internal/posts/domain"
)

type postRouterCtx struct {
	echo        *echo.Echo
	postUsecase domain.PostUsecase
}

func NewPostsRouter(echo *echo.Echo, postUsecase domain.PostUsecase) {
	routerGroup := echo.Group("/posts")
	routerCtx := &postRouterCtx{
		echo,
		postUsecase,
	}

	routerGroup.POST("", routerCtx.createPost)
	routerGroup.GET("/:id", routerCtx.getPost)
	routerGroup.GET("", routerCtx.getPosts)
}

func (ctx *postRouterCtx) createPost(c echo.Context) error {
	var post PostIn

	if err := c.Bind(&post); err != nil {
		return HTTPError{
			Message: "Invalid request body",
		}.ErrUnprocessableEntity()
	}

	post_, err := ctx.postUsecase.Create(c.Request().Context(), post.Title, post.Author, post.Slug, post.Description, post.Content)
	if err != nil {
		return handleErr(err)
	}

	postOut := &PostOut{
		ID:          post_.PublicID,
		Title:       post_.Title,
		Author:      post_.Author,
		Slug:        post_.Slug,
		Status:      post_.Status,
		Description: post_.Description,
		Content:     post_.Content,
		PublishedAt: post_.PublishedAt,
		CreatedAt:   post_.CreatedAt,
		UpdatedAt:   post_.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, postOut)
}

func (ctx *postRouterCtx) getPost(c echo.Context) error {
	var params GetPostParams

	if err := c.Bind(&params); err != nil {
		return HTTPError{
			Message: "Invalid params",
		}.BadRequest()
	}

	post, err := ctx.postUsecase.Get(c.Request().Context(), params.ID)

	if err != nil {
		return handleErr(err)
	}

	postOut := &PostOut{
		ID:          post.PublicID,
		Title:       post.Title,
		Author:      post.Author,
		Slug:        post.Slug,
		Status:      post.Status,
		Description: post.Description,
		Content:     post.Content,
		PublishedAt: post.PublishedAt,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
	}

	return c.JSON(http.StatusOK, postOut)
}

func (ctx *postRouterCtx) getPosts(c echo.Context) error {
	posts, err := ctx.postUsecase.GetPosts(c.Request().Context())

	if err != nil {
		return handleErr(err)
	}

	postsOut := []*PostOut{}

	for _, post := range posts {
		postsOut = append(postsOut, &PostOut{
			ID:          post.PublicID,
			Title:       post.Title,
			Author:      post.Author,
			Slug:        post.Slug,
			Status:      post.Status,
			Description: post.Description,
			Content:     post.Content,
			PublishedAt: post.PublishedAt,
			CreatedAt:   post.CreatedAt,
			UpdatedAt:   post.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, &PostsOut{
		Data: postsOut,
	})
}

func handleErr(err error) error {
	switch err {
	case domain.ErrPostNotFound:
		return HTTPError{
			Message: "Post not found",
		}.NotFound()
	default:
		return HTTPError{
			Message: "Internal server error",
		}.InternalServerError()
	}
}
