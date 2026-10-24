package speedtest

import (
	"context"
	"net/http"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/labstack/echo/v5"
)

type getRepo interface {
	Get(context.Context, uuid.UUID) (*results.Result, error)
}

type workflowChecker interface {
	GetStatus(context.Context, string) (workflow.InstanceStatus, error)
}

type getHandler struct {
	repo   getRepo
	status workflowChecker
}

type GetRequest struct {
	ID uuid.UUID `path:"id" param:"id" description:"The id of the result"`
}

type WorkflowStep struct {
	Name        string              `json:"name"`
	Attempts    int                 `json:"attempts"`
	StartedAt   time.Time           `json:"started_at"`
	CompletedAt time.Time           `json:"completed_at"`
	Status      workflow.StepStatus `json:"status"`
}

type WorkflowStatus struct {
	Version int            `json:"version"`
	Steps   []WorkflowStep `json:"steps"`
}

type GetResponse struct {
	results.Result
	Status *WorkflowStatus `json:"status,omitempty"`
}

func (g *getHandler) Handler() handlers.Handler[GetRequest, GetResponse] {
	return func(c *echo.Context, req GetRequest) (*GetResponse, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Get result")
		defer span.End()

		result, err := g.repo.Get(ctx, req.ID)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		response := &GetResponse{
			Result: *result,
		}

		status, err := g.status.GetStatus(ctx, result.ID.String())
		if err == nil && len(status.Steps) > 0 {
			steps := []WorkflowStep{}
			for _, s := range status.Steps {
				steps = append(steps, WorkflowStep{
					Name:        s.Name,
					Attempts:    s.Attempts,
					Status:      s.Status,
					StartedAt:   s.StartedAt,
					CompletedAt: s.CompletedAt,
				})
			}
			response.Status = &WorkflowStatus{
				Version: status.Version,
				Steps:   steps,
			}
		}

		return response, nil
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
