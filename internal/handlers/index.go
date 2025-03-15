package handlers

import (
	"context"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/akithepriest/chisai.click/internal/errs"
	"github.com/akithepriest/chisai.click/internal/renderer"
	"github.com/akithepriest/chisai.click/internal/store"
	"github.com/akithepriest/chisai.click/internal/views"
	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
)

type apiResponse struct {
	Success    bool        `json:"success"`
	StatusText string      `json:"status_text"`
	D          interface{} `json:"d"`
}

type IndexHandler struct {
	linkStore store.LinkStore
}

func NewIndexHandler(db *mongo.Database) *IndexHandler {
	return &IndexHandler{
		linkStore: store.NewLinkDBStore(db),
	}
}

func (h *IndexHandler) DefineRoutes(e *echo.Group) {
	e.GET("/", h.handleGET)

	e.GET("/:keyword", h.handleGETWithKeyword)
	e.POST("/api/create", h.handlePOSTLink)
}

func (h *IndexHandler) handleGET(c echo.Context) error {
	return renderer.Render(c, http.StatusOK, views.Index())
}

func (h *IndexHandler) handleGETWithKeyword(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	keyword := c.Param("keyword")

	link, err := h.linkStore.GetByKeyword(ctx, keyword)
	if err != nil {
		if err == errs.ErrNotFound {
			return renderer.RenderError(c, http.StatusNotFound, err)
		}
		return renderer.RenderError(c, http.StatusInternalServerError, err)
	}
	return c.Redirect(http.StatusSeeOther, link.RedirectLink)
}

func (h *IndexHandler) handlePOSTLink(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	keyword, link := c.FormValue("keyword"), c.FormValue("link")
	if keyword == "" || link == "" {
		return c.JSON(http.StatusBadRequest, apiResponse{
			StatusText: http.StatusText(http.StatusBadRequest),
			D:          "Keyword or link not provided",
		})
	}

	validKeyword := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !validKeyword.MatchString(keyword) {
		return c.JSON(http.StatusBadRequest, apiResponse{
			StatusText: http.StatusText(http.StatusBadRequest),
			D:          "Invalid keyword format. Only alphanumeric characters, hyphens, and underscores are allowed.",
		})
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		return c.JSON(http.StatusBadRequest, apiResponse{
			StatusText: http.StatusText(http.StatusBadRequest),
			D:          "Invalid URL format",
		})
	}

	newLink, err := h.linkStore.Create(ctx, keyword, link)
	if err != nil {
		if err == errs.ErrAlreadyExists {
			return c.JSON(http.StatusBadRequest, apiResponse{
				StatusText: http.StatusText(http.StatusBadRequest),
				D:          "Keyword already claimed",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, apiResponse{
				StatusText: http.StatusText(http.StatusInternalServerError),
				D:          err.Error(),
			})
		}
	}
	return c.JSON(http.StatusOK, apiResponse{
		Success:    true,
		StatusText: http.StatusText(http.StatusOK),
		D:          newLink,
	})
}
