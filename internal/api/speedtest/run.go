package speedtest

import (
	"context"
	"net/http"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

type actorInvoker interface {
	Invoke(
		context.Context,
		string,
		string,
		string,
		any,
		...actor.InvokeOption,
	) (actor.Envelope, error)
}

type rateLimiter interface {
	Allow(context.Context, string) (bool, time.Duration, error)
}

type runHandler struct {
	service actorInvoker
	limiter rateLimiter
}

type RunResponse struct {
	ID uuid.UUID `json:"id"`
}

func (r *runHandler) Handler() handlers.Handler[any, RunResponse] {
	return func(c *echo.Context, req any) (*RunResponse, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Run speedtest")
		defer span.End()

		limitCtx, span := tracing.NewSpan(ctx, "Check rate limiter")
		defer span.End()
		if allowed, _, err := r.limiter.Allow(limitCtx, "test"); err != nil || !allowed {
			return nil, echo.ErrTooManyRequests
		}
		span.End()

		id, _ := uuid.New()

		runCtx, span := tracing.NewSpan(ctx, "Invoke actor run")
		defer span.End()
		raw, err := r.service.Invoke(runCtx, ookla.Actor, id.String(), ookla.Onboard, nil)
		if err != nil {
			return nil, httperrors.Stack(err)
		}
		span.End()

		var result results.Result
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
