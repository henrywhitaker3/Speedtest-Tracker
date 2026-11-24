package config

import (
	"net/http"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/henrywhitaker3/windowframe/v2/francis"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

type listServersHandler struct {
	invoker francis.Invoker
}

func (l *listServersHandler) Handler() handlers.Handler[any, []ookla.Server] {
	return func(c *echo.Context, req any) (*[]ookla.Server, error) {
		ctx, span := tracing.NewSpan(c.Request().Context(), "List servers")
		defer span.End()

		data, err := l.invoker.Invoke(
			ctx,
			ookla.Actor,
			actor.SingletonActorID,
			ookla.ListServers,
			nil,
		)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		var servers []ookla.Server
		if err := data.Decode(&servers); err != nil {
			return nil, httperrors.Stack(err)
		}

		return &servers, nil
	}
}

func (l *listServersHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *listServersHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "List servers",
		Description:  "List nearby Speedtest servers",
		ID:           new("listServers"),
		Tag:          "Config",
		Path:         "/v1/config/servers",
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
