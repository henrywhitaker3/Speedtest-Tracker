package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"slices"
	"strings"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/windowframe/v2/tracing"
	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel/attribute"
)

const (
	sessionCookieName = "__Host-speedtest-session"
)

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			_, span := tracing.NewSpan(c.Request().Context(), "Create session")
			defer span.End()

			if header := c.Request().Header.Get("Authorization"); header != "" {
				return next(c)
			}

			if cookie, err := c.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
				return next(c)
			}

			sessionID, err := newSessionID()
			if err != nil {
				return fmt.Errorf("generate session ID: %w", err)
			}

			setSessionCookie(c, sessionID)

			span.End()

			return next(c)
		}
	}
}

type tokenRepo interface {
	Validate(context.Context, string) error
}

type AuthMiddlewareOpts struct {
	// Routes are a slice of routes that should be protected behind auth
	Routes []string

	Service   actorInvoker
	TokenRepo tokenRepo
}

func AuthMiddleware(opts AuthMiddlewareOpts) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx, span := tracing.NewSpan(c.Request().Context(), "Authentication")
			defer span.End()

			if !slices.Contains(opts.Routes, c.Path()) {
				return next(c)
			}

			authHeader := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
			if authHeader != "" {
				span.SetAttributes(attribute.String("auth.method", "api"))
				if err := handleTokenAuth(ctx, authHeader, opts.TokenRepo); err != nil {
					return err
				}
			} else {
				span.SetAttributes(attribute.String("auth.method", "session"))
				if err := handleSessionAuth(ctx, c, opts.Service); err != nil {
					return err
				}
			}
			span.End()

			return next(c)
		}
	}
}

func handleTokenAuth(ctx context.Context, token string, repo tokenRepo) error {
	if err := repo.Validate(ctx, token); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	return nil
}

func handleSessionAuth(ctx context.Context, c *echo.Context, service actorInvoker) error {
	sessionID, err := cookieValue(c, sessionCookieName)
	if err != nil {
		return err
	}

	if _, err := service.Invoke(
		ctx,
		authactor.Actor,
		sessionID,
		authactor.Session,
		nil,
	); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
	}

	return nil
}

func newSessionID() (string, error) {
	const bytes = 32
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func setSessionCookie(c *echo.Context, value string) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}

func cookieValue(c *echo.Context, name string) (string, error) {
	cookie, err := c.Cookie(name)
	if err != nil || cookie.Value == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "missing session cookie")
	}
	return cookie.Value, nil
}
