package auth

import (
	"fmt"
	"net/http"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/labstack/echo/v5"
)

type sessionHandler struct {
	service actorInvoker
}

func (s *sessionHandler) Handler() handlers.Handler[any, authactor.AuthenticateResult] {
	return func(ctx *echo.Context, _ any) (*authactor.AuthenticateResult, error) {
		sessionID, err := cookieValue(ctx, sessionCookieName)
		if err != nil {
			return nil, err
		}

		raw, err := s.service.Invoke(
			ctx.Request().Context(),
			authactor.Actor,
			sessionID,
			authactor.Session,
			nil,
		)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
		}

		var response authactor.AuthenticateResult
		if err := raw.Decode(&response); err != nil {
			return nil, fmt.Errorf("decode session response: %w", err)
		}

		return &response, nil
	}
}

func (s *sessionHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (s *sessionHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Auth session",
		Tag:          "Auth",
		Path:         "/v1/auth/session",
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		GenerateSpec: false,
	}
}
