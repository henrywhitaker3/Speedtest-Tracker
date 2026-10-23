package speedtest

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type purgeRepo interface {
	PurgeErrors(context.Context) error
}

type purgeHandler struct {
	repo purgeRepo
}

func (p *purgeHandler) Handler() handlers.Handler[any, any] {
	return func(c *echo.Context, req any) (*any, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Purge errored results")
		defer span.End()

		if err := p.repo.PurgeErrors(ctx); err != nil {
			return nil, httperrors.Stack(err)
		}

		return nil, nil
	}
}

func (p *purgeHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (p *purgeHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Purge errors",
		Description:  "Delete test results that have an error",
		Tag:          "Speedtest",
		Path:         "/v1/speedtest/purge",
		Method:       http.MethodDelete,
		Code:         http.StatusOK,
		GenerateSpec: true,
	}
}
