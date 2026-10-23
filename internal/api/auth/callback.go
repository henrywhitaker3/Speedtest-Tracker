package auth

import (
	"context"
	"fmt"
	"net/http"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

type callbackHandler struct {
	service            actorInvoker
	successRedirectURL string
}

type actorInvoker interface {
	Invoke(
		context.Context,
		string,
		string,
		string,
		any,
		...actor.InvokeOption,
	) (actor.Envelope, error)
}

// CallbackRequest contains the authorization code and CSRF state returned by
// the OIDC issuer. The auth actor validates and exchanges the code.
type CallbackRequest struct {
	Code  string `query:"code"`
	State string `query:"state"`
	Error string `query:"error"`
}

func (c *callbackHandler) Handler() handlers.Handler[CallbackRequest, any] {
	return func(ctx *echo.Context, req CallbackRequest) (*any, error) {
		if req.Error != "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "OIDC issuer returned "+req.Error)
		}
		if req.Code == "" || req.State == "" {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "missing OIDC code or state")
		}

		sessionID, err := cookieValue(ctx, sessionCookieName)
		if err != nil {
			return nil, err
		}
		_, err = c.service.Invoke(
			ctx.Request().Context(),
			authactor.Actor,
			sessionID,
			authactor.CompleteLogin,
			authactor.CompleteLoginRequest{Code: req.Code, State: req.State},
		)
		if err != nil {
			return nil, fmt.Errorf("authenticate session: %w", err)
		}

		if err := ctx.Redirect(http.StatusSeeOther, c.successRedirectURL); err != nil {
			return nil, fmt.Errorf("redirect to dashboard: %w", err)
		}
		return nil, nil
	}
}

func (c *callbackHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (c *callbackHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Auth callback",
		Tag:          "Auth",
		Path:         "/v1/auth/callback",
		Code:         http.StatusSeeOther,
		Method:       http.MethodGet,
		GenerateSpec: false,
	}
}
