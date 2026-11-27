package config

import (
	"context"
	"net/http"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/henrywhitaker3/windowframe/v2/francis"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

type jobLister interface {
	ListJobs(ctx context.Context, actorType string, actorID string) ([]actor.JobInfo, error)
}

type historicalWeatherStatusHandler struct {
	job    francis.Dispatcher
	lister jobLister
}

func (h *historicalWeatherStatusHandler) Handler() handlers.Handler[any, actor.JobInfo] {
	return func(c *echo.Context, req any) (*actor.JobInfo, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Get weather backfill status")
		defer span.End()

		jobs, err := h.lister.ListJobs(ctx, weather.Actor, actor.SingletonActorID)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		for _, j := range jobs {
			if j.Method == weather.Backfill {
				return &j, nil
			}
		}

		return nil, echo.ErrNotFound
	}
}

func (h *historicalWeatherStatusHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (h *historicalWeatherStatusHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Historical Weather Data Status",
		ID:           new("historicalWeatherStatus"),
		Code:         http.StatusOK,
		Description:  "Get the status of historical weather data backfill",
		Method:       http.MethodGet,
		Path:         "/v1/config/weather/backfill",
		GenerateSpec: true,
		Tag:          "Config",
		Auths: []handlers.Auth{
			{
				Enabled: true,
				Name:    authapi.SessionSecurityScheme,
			},
		},
	}
}
