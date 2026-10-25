package speedtest

import (
	"context"
	"net/http"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/duration"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type timeseriesRepo interface {
	Timeseries(context.Context, time.Duration) ([]results.Stat, error)
}

type timeseriesHandler struct {
	repo timeseriesRepo
}

type TimeseriesRequest struct {
	Period duration.StringDuration `query:"period" description:"The time range to pick from" validate:"timeseries_period" default:"24h"`
}

func (t *timeseriesHandler) Handler() handlers.Handler[TimeseriesRequest, []results.Stat] {
	return func(c *echo.Context, req TimeseriesRequest) (*[]results.Stat, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Speedtest timeseries")
		defer span.End()

		if req.Period == 0 {
			req.Period = duration.StringDuration(time.Hour * 24)
		}

		stats, err := t.repo.Timeseries(ctx, req.Period.Duration())
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &stats, nil
	}
}

func (t *timeseriesHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (t *timeseriesHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Timeseries",
		Description:  "Get timeseries data for speedtest results",
		ID:           new("timeseries"),
		Tag:          "Speedtest",
		Method:       http.MethodGet,
		Code:         http.StatusOK,
		Path:         "/v1/speedtest/timeseries",
		GenerateSpec: true,
	}
}
