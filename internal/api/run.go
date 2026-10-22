package api

import (
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/ookla"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

type runHandler struct {
	service *actor.Service
}

type RunResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r *runHandler) Handler() handlers.Handler[any, RunResponse] {
	return func(c *echo.Context, req any) (*RunResponse, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Run speedtest")
		defer span.End()

		id, _ := uuid.New()

		raw, err := r.service.Invoke(ctx, ookla.Actor, id.String(), ookla.Onboard, nil)
		if err != nil {
			return nil, httperrors.Stack(err)
		}
		var result ookla.Result
		if err := raw.Decode(&result); err != nil {
			return nil, httperrors.Stack(err)
		}

		return &RunResponse{
			ID: result.ID,
		}, nil
	}
}

func (r *runHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (r *runHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Run speedtest",
		Tag:          "Speedtest",
		Method:       http.MethodPost,
		Path:         "/v1/speedtest/run",
		Code:         http.StatusOK,
		Description:  "Run a new speedtest",
		GenerateSpec: true,
	}
}
