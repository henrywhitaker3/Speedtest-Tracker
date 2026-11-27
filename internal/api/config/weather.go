package config

import (
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/labstack/echo/v5"
)

type getCurrentWeatherHandler struct {
	p weather.Provider
}

func (h *getCurrentWeatherHandler) Handler() handlers.Handler[any, weather.WeatherData] {
	return func(c *echo.Context, _ any) (*weather.WeatherData, error) {
		if h.p == nil {
			return nil, echo.NewHTTPError(http.StatusServiceUnavailable, "weather is not configured")
		}
		data, err := h.p.Get(c.Request().Context())
		if err != nil {
			return nil, httperrors.Stack(err)
		}
		return data, nil
	}
}

func (h *getCurrentWeatherHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (h *getCurrentWeatherHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Current weather",
		Description:  "Get current weather for the configured location",
		ID:           new("getCurrentWeather"),
		Tag:          "Config",
		Path:         "/v1/config/weather",
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		GenerateSpec: true,
	}
}
