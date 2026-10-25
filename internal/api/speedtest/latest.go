package speedtest

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type latestRepo interface {
	Latest(context.Context) (*results.Result, error)
}

type latestHandler struct {
	repo latestRepo
}

func (l *latestHandler) Handler() handlers.Handler[any, Result] {
	return func(c *echo.Context, req any) (*Result, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Get latest result")
		defer span.End()

		latest, err := l.repo.Latest(ctx)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return new(convert(c, *latest)), nil
	}
}

func (l *latestHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *latestHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Latest",
		Description:  "Get the latest speedtest result",
		ID:           new("latestResult"),
		Path:         "/v1/speedtest/latest",
		Method:       http.MethodGet,
		Code:         http.StatusOK,
		Tag:          "Speedtest",
		GenerateSpec: true,
	}
}
