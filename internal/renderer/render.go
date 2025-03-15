package renderer

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/akithepriest/chisai.click/internal/views"
	"github.com/labstack/echo"
)

// Set the http status and render specified templ page.
func Render(c echo.Context, status int, cmp templ.Component) error {
	c.Response().WriteHeader(status)

	return cmp.Render(c.Request().Context(), c.Response())
}

// Render error templ.
func RenderError(c echo.Context, status int, err error) error {
	c.Response().WriteHeader(status)

	return views.Error(status, http.StatusText(status)).Render(c.Request().Context(), c.Response())
}