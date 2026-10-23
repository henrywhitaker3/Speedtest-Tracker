package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/henrywhitaker3/windowframe/v2/http/httperrors"
	"github.com/henrywhitaker3/windowframe/v2/http/validation"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
)

type createRepo interface {
	Create(context.Context, string) (string, auth.Token, error)
}

type createHandler struct {
	repo createRepo
}

type CreateTokenRequest struct {
	Name string `json:"name" description:"The name to identify the token by" validate:"required,min=5,max=255"`
}

type CreateTokenResponse struct {
	Token  auth.Token `json:"token"`
	Secret string     `json:"secret"`
}

func (h *createHandler) Handler() handlers.Handler[CreateTokenRequest, CreateTokenResponse] {
	return func(c *echo.Context, req CreateTokenRequest) (*CreateTokenResponse, error) {
		name := strings.TrimSpace(req.Name)
		if name == "" {
			return nil, validation.Build().With("name", "required")
		}

		ctx, span := tracing.NewSpan(c.Request().Context(), "Create auth token")
		defer span.End()

		secret, token, err := h.repo.Create(ctx, name)
		if err != nil {
			return nil, httperrors.Stack(err)
		}

		return &CreateTokenResponse{
			Token:  token,
			Secret: secret,
		}, nil
	}
}

func (h *createHandler) Middleware() []echo.MiddlewareFunc {
	// TODO: use admin mw
	return nil
}

func (g *createHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Create token",
		Tag:          "Auth",
		Path:         "/v1/auth/tokens",
		Method:       http.MethodPost,
		Code:         http.StatusOK,
		GenerateSpec: true,
	}
}
