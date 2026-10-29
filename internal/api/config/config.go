// Package config
package config

import (
	"net/http"

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

type ConfigResponse struct {
	Auth     AuthResponse `json:"auth"`
	Schedule string       `json:"schedule"`
	ServerID []int        `json:"server_id"`
}

func (h *configHandler) Handler() handlers.Handler[any, ConfigResponse] {
	return func(c *echo.Context, req any) (*ConfigResponse, error) {
		return &ConfigResponse{
			Auth: AuthResponse{
				Enabled:    h.conf.Auth.Enabled,
				AdminGroup: h.conf.Auth.AdminGroup,
			},
			Schedule: h.conf.Schedule,
			ServerID: h.conf.Server,
		}, nil
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
