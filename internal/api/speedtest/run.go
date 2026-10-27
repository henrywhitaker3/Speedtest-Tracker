package speedtest

import (
	"context"
	"net/http"
	"time"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/labstack/echo/v5"
)

type resultCreator interface {
	Create(context.Context) (*results.Result, error)
	Delete(context.Context, uuid.UUID) error
}

type worflowRunner interface {
	Start(context.Context, any, ...workflow.StartOption) (string, bool, error)
}

type rateLimiter interface {
	Allow(context.Context, string) (bool, time.Duration, error)
}

type runHandler struct {
	repo     resultCreator
	workflow worflowRunner
	limiter  rateLimiter
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

		runCtx, span := tracing.NewSpan(ctx, "Start speedtest workflow")
		defer span.End()

		result, err := r.repo.Create(runCtx)
		if err != nil {
			return nil, httperrors.Stack(err)
		}
		if _, _, err := r.workflow.Start(
			runCtx,
			ookla.RunWorkflow{
				Result: *result,
			},
			workflow.WithInstanceID(result.ID.String()),
		); err != nil {
			_ = r.repo.Delete(ctx, result.ID)
			return nil, httperrors.Stack(err)
		}
		span.End()

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
		ID:           new("run"),
		Method:       http.MethodPost,
		Path:         "/v1/speedtest/run",
		Code:         http.StatusOK,
		Description:  "Run a new speedtest",
		GenerateSpec: true,
		Auths: []handlers.Auth{
			{
				Enabled: true,
				Name:    authapi.APITokenSecurityScheme,
			},
			{
				Enabled: true,
				Name:    authapi.SessionSecurityScheme,
			},
		},
	}
}
