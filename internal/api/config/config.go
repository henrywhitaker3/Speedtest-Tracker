// Package config
package config

import (
	"net/http"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/labstack/echo/v5"
)

type configHandler struct {
	conf config.Config
}

type AuthResponse struct {
	Enabled    bool   `json:"enabled"`
	AdminGroup string `json:"admin_group"`
}

type WeatherResponse struct {
	Enabled   bool    `json:"enabled"`
	Provider  string  `json:"provider"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type ConfigResponse struct {
	Auth     AuthResponse    `json:"auth"`
	Schedule string          `json:"schedule"`
	ServerID []int           `json:"server_id"`
	Weather  WeatherResponse `json:"weather"`
}

func (h *configHandler) Handler() handlers.Handler[any, ConfigResponse] {
	return func(c *echo.Context, req any) (*ConfigResponse, error) {
		resp := &ConfigResponse{
			Auth: AuthResponse{
				Enabled:    h.conf.Auth.Enabled,
				AdminGroup: h.conf.Auth.AdminGroup,
			},
			Schedule: h.conf.Schedule,
			ServerID: h.conf.Server,
			Weather: WeatherResponse{
				Enabled:   h.conf.Weather.Enabled,
				Provider:  h.conf.Weather.Provider,
				Longitude: h.conf.Weather.Longitude,
				Latitude:  h.conf.Weather.Latitude,
			},
		}

		if !auth.IsAuthenticated(c) {
			resp.ServerID = nil
		}

		return resp, nil
	}
}

func (h *configHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (h *configHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Config",
		Tag:          "Config",
		ID:           new("getConfig"),
		Path:         "/v1/config",
		Method:       http.MethodGet,
		Code:         http.StatusOK,
		Description:  "Get the app config",
		GenerateSpec: true,
	}
}
