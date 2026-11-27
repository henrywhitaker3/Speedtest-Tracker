package config

import (
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

type historicalWeatherBackfillHandler struct {
	dispatcher francis.Dispatcher
}

type HistoricalWeatherBackfillResponse struct {
	JobID   string `json:"job_id"`
	Created bool   `json:"created"`
}

func (h *historicalWeatherBackfillHandler) Handler() handlers.Handler[any, HistoricalWeatherBackfillResponse] {
	return func(c *echo.Context, _ any) (*HistoricalWeatherBackfillResponse, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "Start weather data backfill")
		defer span.End()

		id, created, err := h.dispatcher.Dispatch(
			ctx,
			weather.Actor,
			actor.SingletonActorID,
			weather.Backfill,
			nil,
			actor.WithIdempotencyKey(weather.Backfill),
		)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &HistoricalWeatherBackfillResponse{JobID: id, Created: created}, nil
	}
}

func (h *historicalWeatherBackfillHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (h *historicalWeatherBackfillHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Historical Weather Data Backfill",
		ID:           new("historicalWeatherBackfill"),
		Code:         http.StatusOK,
		Description:  "Backfill historical weather data",
		Method:       http.MethodPost,
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
