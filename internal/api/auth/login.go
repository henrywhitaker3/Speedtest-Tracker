package auth

import (
	"fmt"
	"net/http"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/labstack/echo/v5"
)

type loginHandler struct {
	service actorInvoker
}

func (l *loginHandler) Handler() handlers.Handler[any, any] {
	return func(c *echo.Context, _ any) (*any, error) {
		sessionID, err := cookieValue(c, SessionCookieName)
		if err != nil {
			return nil, err
		}

		raw, err := l.service.Invoke(
			c.Request().Context(),
			authactor.Actor,
			sessionID,
			authactor.BeginLogin,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("start login transaction: %w", err)
		}
		var response authactor.LoginResponse
		if err := raw.Decode(&response); err != nil {
			return nil, fmt.Errorf("decode login response: %w", err)
		}
		if response.RedirectURL == "" {
			return nil, fmt.Errorf("decode login response: redirect URL is required")
		}

		if err := c.Redirect(http.StatusFound, response.RedirectURL); err != nil {
			return nil, fmt.Errorf("redirect to OIDC issuer: %w", err)
		}
		return nil, nil
	}
}

func (l *loginHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *loginHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Auth login",
		Tag:          "Auth",
		Path:         "/v1/auth/login",
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		GenerateSpec: false,
	}
}
