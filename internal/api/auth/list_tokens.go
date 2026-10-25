package auth

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type listRepo interface {
	List(context.Context) ([]auth.Token, error)
}

type listTokensHandler struct {
	repo listRepo
}

func (l *listTokensHandler) Handler() handlers.Handler[any, []auth.Token] {
	return func(c *echo.Context, req any) (*[]auth.Token, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "List auth tokens")
		defer span.End()

		tokens, err := l.repo.List(ctx)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &tokens, nil
	}
}

func (l *listTokensHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *listTokensHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "List tokens",
		Tag:          "Auth",
		ID:           new("listTokens"),
		Path:         "/v1/auth/tokens",
		Method:       http.MethodGet,
		Code:         http.StatusOK,
		GenerateSpec: true,
	}
}
