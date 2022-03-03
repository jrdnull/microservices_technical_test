package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

type articleHandler struct {
	articles *internal.ArticleService
	users    *internal.UserService
	cfg      Config
}

func (h articleHandler) mount(g *echo.Group) {
	g.GET("/articles", h.listArticles)
	g.GET("/articles/feed", h.getUserFeed)
}

func (h articleHandler) listArticles(c echo.Context) error {
	var opts internal.ListOptions
	opts.Page, _ = strconv.Atoi(c.QueryParam("page"))
	var tagFilter string
	if v := c.QueryParam("all_tags"); v != "" {
		tagFilter = v
		opts.MatchAllTags = true
	} else if v := c.QueryParam("any_tags"); v != "" {
		tagFilter = v
	}
	if tagFilter != "" {
		ids, err := parseIDs(tagFilter)
		if err != nil {
			return fmt.Errorf("tag filter: %w", err)
		}
		opts.TagIDs = ids
	}

	if opts.MatchAllTags && len(opts.TagIDs) != h.cfg.ArticleTagFilterInputs {
		// enforce predicate required by design
		return internal.NewValidationErrorf("all_tags requires exactly %d ids", h.cfg.ArticleTagFilterInputs)
	}

	articles, err := h.articles.List(c.Request().Context(), opts)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles": articles,
	})
}

func (h articleHandler) getUserFeed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))

	ctx := c.Request().Context()
	user, err := h.users.Get(ctx, mustGetUserID(ctx))
	if err != nil {
		return fmt.Errorf("get authenticated user: %w", err)
	}
	articles, err := h.articles.List(ctx, internal.ListOptions{
		TagIDs: user.NewsFeedTags,
		Page:   page,
	})
	if err != nil {
		return fmt.Errorf("list articles: %w", err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"articles": articles,
	})
}
