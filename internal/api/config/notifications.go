package config

import (
	"net/http"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/labstack/echo/v5"
)

type notificationsHandler struct {
	conf config.Notifications
}

func (n *notificationsHandler) Handler() handlers.Handler[any, config.Notifications] {
	return func(c *echo.Context, req any) (*config.Notifications, error) {
		return &n.conf, nil
	}
}

func (n *notificationsHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (n *notificationsHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Notifications",
		Description:  "Show configured notification channels",
		ID:           new("getNotifications"),
		Tag:          "Config",
		Path:         "/v1/config/notifications",
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		GenerateSpec: true,
		Auths: []handlers.Auth{
			{
				Enabled: true,
				Name:    authapi.SessionSecurityScheme,
			},
		},
	}
}
