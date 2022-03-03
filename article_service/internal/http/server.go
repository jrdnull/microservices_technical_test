package http

import (
	"context"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/jrdnull/microservices_technical_test/article_service/internal"
)

// Config provides configuration to the handlers.
type Config struct {
	// Specifies the required amount of tag ids to provide to the
	// /articles tag filter.
	ArticleTagFilterInputs int
}

// NewServer returns a new Echo server with resources mounted.
func NewServer(
	articleSrv *internal.ArticleService, userSrv *internal.UserService, cfg Config,
) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = errorHandler
	e.Use(middleware.Recover(), otelecho.Middleware(internal.ServiceName), middleware.Logger())

	// TODO: real authentication
	articles := e.Group("/v1", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:auth_key",
		Validator: func(key string, c echo.Context) (bool, error) {
			var userID int64
			switch key {
			case "user1":
				userID = 1
			case "user2":
				userID = 2
			default:
				return false, internal.ErrUnauthorized
			}
			c.SetRequest(c.Request().WithContext(withUserID(c.Request().Context(), userID)))
			return true, nil
		},
		ErrorHandler: func(err error, c echo.Context) error {
			return internal.ErrUnauthorized
		},
	}))
	articleHandler{articleSrv, userSrv, cfg}.mount(articles)

	return e
}

func parseIDs(s string) ([]int64, error) {
	var ids []int64
	for _, sid := range strings.Split(s, ",") {
		id, _ := strconv.ParseInt(strings.TrimSpace(sid), 10, 64)
		if id < 1 {
			return nil, internal.NewValidationError("invalid id")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

type userIDKey struct{}

func withUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, userIDKey{}, id)
}

func mustGetUserID(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDKey{}).(int64)
	if !ok {
		panic("userID not set on context, forgot auth middleware?")
	}
	return userID
}
