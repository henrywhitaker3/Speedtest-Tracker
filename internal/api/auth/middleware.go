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
	SessionCookieName      = "__Host-speedtest-session"
	APITokenSecurityScheme = "api_token"
	SessionSecurityScheme  = "session"

	authEnabledContextKey   = "auth-enabled"
	authenticatedContextKey = "authenticated"
	adminContextKey         = "admin"
)

func AddContext(authEnabled bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if authEnabled {
				c.Set(authEnabledContextKey, true)
			} else {
				c.Set(authEnabledContextKey, false)
			}
			return next(c)
		}
	}
}

func IsAdmin(c *echo.Context) bool {
	if !authEnabled(c) {
		return true
	}

	if !IsAuthenticated(c) {
		return false
	}

	// If auth is enabled (or this is a context assembled by a focused handler
	// test), check whether authentication established administrator access.
	if admin, ok := c.Get(adminContextKey).(bool); admin && ok {
		return true
	}

	return false
}

func authEnabled(c *echo.Context) bool {
	authEnabled, configured := c.Get(authEnabledContextKey).(bool)
	if !configured {
		return false
	}
	return authEnabled
}

func IsAuthenticated(c *echo.Context) bool {
	if !authEnabled(c) {
		return true
	}

	authenticated, ok := c.Get(authenticatedContextKey).(bool)
	return authenticated && ok
}

func SessionMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			_, span := tracing.NewSpan(c.Request().Context(), "Create session")
			defer span.End()

			if header := c.Request().Header.Get("Authorization"); header != "" {
				return next(c)
			}

			if cookie, err := c.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
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

type ProtectedRoute struct {
	Method string
	Path   string
}

type ProtectedRoutes []ProtectedRoute

func (p ProtectedRoutes) Contains(c *echo.Context) bool {
	for _, r := range p {
		if c.Request().Method == r.Method && c.Path() == r.Path {
			return true
		}
	}
	return false
}

type AuthMiddlewareOpts struct {
	// Routes are a slice of routes that MUST be authenticated
	Routes ProtectedRoutes

	Service    actorInvoker
	TokenRepo  tokenRepo
	AdminGroup string
}

// MiddlewareOpts configures the API authentication middleware stack.
type MiddlewareOpts struct {
	Enabled bool

	Service     actorInvoker
	TokenRepo   tokenRepo
	AdminGroup  string
	Routes      ProtectedRoutes
	AdminRoutes ProtectedRoutes
}

// Middleware returns the authentication middleware in the order used by the
// API. Authentication context is always available; request authentication and
// administrator authorisation are only needed when auth is enabled.
func Middleware(opts MiddlewareOpts) []echo.MiddlewareFunc {
	middleware := []echo.MiddlewareFunc{AddContext(opts.Enabled)}
	if !opts.Enabled {
		return middleware
	}

	return append(middleware,
		SessionMiddleware(),
		AuthMiddleware(AuthMiddlewareOpts{
			Service:    opts.Service,
			TokenRepo:  opts.TokenRepo,
			AdminGroup: opts.AdminGroup,
			Routes:     opts.Routes,
		}),
		Admin(AdminOpts{Routes: opts.AdminRoutes}),
	)
}

func AuthMiddleware(opts AuthMiddlewareOpts) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			ctx, span := tracing.NewSpan(c.Request().Context(), "Authenticate")
			defer span.End()

			var err error
			authHeader := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
			if authHeader != "" {
				span.SetAttributes(attribute.String("auth.method", "api"))
				err = handleTokenAuth(ctx, authHeader, opts.TokenRepo)
			} else {
				span.SetAttributes(attribute.String("auth.method", "session"))
				err = handleSessionAuth(ctx, c, opts.Service, opts.AdminGroup)
			}

			if err == nil {
				c.Set(authenticatedContextKey, true)
			}

			if !opts.Routes.Contains(c) {
				span.End()
				return next(c)
			}

			if err != nil {
				return err
			}

			span.End()
			return next(c)
		}
	}
}

type AdminOpts struct {
	Routes ProtectedRoutes
}

func Admin(opts AdminOpts) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			_, span := tracing.NewSpan(c.Request().Context(), "Authorise")
			defer span.End()

			if !opts.Routes.Contains(c) {
				span.End()
				return next(c)
			}

			if IsAdmin(c) {
				span.End()
				return next(c)
			}

			return echo.NewHTTPError(http.StatusForbidden, "forbidden")
		}
	}
}

func handleTokenAuth(ctx context.Context, token string, repo tokenRepo) error {
	if err := repo.Validate(ctx, token); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
	}
	return nil
}

func handleSessionAuth(
	ctx context.Context,
	c *echo.Context,
	service actorInvoker,
	adminGroup string,
) error {
	sessionID, err := cookieValue(c, SessionCookieName)
	if err != nil {
		return err
	}

	sessionRaw, err := service.Invoke(
		ctx,
		authactor.Actor,
		sessionID,
		authactor.Session,
		nil,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired session")
	}

	if sessionRaw != nil {
		var session authactor.AuthenticateResult
		if err := sessionRaw.Decode(&session); err != nil {
			return fmt.Errorf("decode auth session: %w", err)
		}
		if slices.Contains(session.Groups, adminGroup) {
			c.Set(adminContextKey, true)
		}
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
		Name:     SessionCookieName,
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
