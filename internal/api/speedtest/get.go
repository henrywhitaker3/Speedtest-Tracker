package speedtest

import (
	"context"
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/labstack/echo/v5"
)

type getRepo interface {
	Get(context.Context, uuid.UUID) (*results.Result, error)
}

type getHandler struct {
	repo getRepo
}

type GetRequest struct {
	ID uuid.UUID `path:"id" param:"id" description:"The id of the result"`
}

func (g *getHandler) Handler() handlers.Handler[GetRequest, results.Result] {
	return func(c *echo.Context, req GetRequest) (*results.Result, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Get result")
		defer span.End()

		result, err := g.repo.Get(ctx, req.ID)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return result, nil
	}
}

func (g *getHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (g *getHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Get result",
		Tag:          "Speedtest",
		Path:         "/v1/speedtest/:id",
		Method:       http.MethodGet,
		Code:         http.StatusOK,
		Description:  "Get a speedtest result",
		GenerateSpec: true,
	}
}
