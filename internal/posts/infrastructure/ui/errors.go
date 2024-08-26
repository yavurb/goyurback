package ui

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HTTPError struct {
	Message string `json:"message"`
}

func (e HTTPError) InternalServerError() error {
	err := echo.ErrInternalServerError
	err.Message = e.Message

	return err
}

func (e HTTPError) BadRequest() error {
	return echo.NewHTTPError(http.StatusBadRequest, e.Message)
}

func (e HTTPError) NotFound() error {
	err := echo.ErrNotFound
	err.Message = e.Message

	return err
}

func (e HTTPError) Unauthorized() error {
	return echo.NewHTTPError(http.StatusUnauthorized, e.Message)
}

func (e HTTPError) Forbidden() error {
	return echo.NewHTTPError(http.StatusForbidden, e.Message)
}

func (e HTTPError) Conflict() error {
	return echo.NewHTTPError(http.StatusConflict, e.Message)
}

func (e HTTPError) ErrUnprocessableEntity() error {
	err := echo.ErrUnprocessableEntity

	err.Message = e.Message

	return err
}
