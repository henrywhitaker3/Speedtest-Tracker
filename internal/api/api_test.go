package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

func TestAuthMiddlewareRouteConfiguration(t *testing.T) {
	tests := []struct {
		method    string
		path      string
		auth      bool
		adminOnly bool
	}{
		{method: http.MethodPost, path: "/v1/speedtest/run", auth: true},
		{method: http.MethodDelete, path: "/v1/speedtest/purge", auth: true, adminOnly: true},
		{method: http.MethodDelete, path: "/v1/speedtest/:id", auth: true, adminOnly: true},
		{method: http.MethodGet, path: "/v1/auth/tokens", auth: true, adminOnly: true},
		{method: http.MethodPost, path: "/v1/auth/tokens", auth: true, adminOnly: true},
		{method: http.MethodDelete, path: "/v1/auth/tokens/:id", auth: true, adminOnly: true},
		{method: http.MethodGet, path: "/v1/config/notifications", auth: true, adminOnly: true},
		{method: http.MethodGet, path: "/v1/speedtest/latest"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			assert.Equal(t, tt.auth, hasProtectedRoute(authRoutes(), tt.method, tt.path))
			assert.Equal(
				t,
				tt.adminOnly,
				hasProtectedRoute(adminRoutes(), tt.method, tt.path),
			)
		})
	}
}

func TestNotificationsEndpointAccess(t *testing.T) {
	tests := []struct {
		name        string
		authEnabled bool
		withCookie  bool
		groups      []string
		wantStatus  int
		wantCalled  bool
	}{
		{
			name:        "rejects an unauthenticated request when auth is enabled",
			authEnabled: true,
			wantStatus:  http.StatusUnauthorized,
		},
		{
			name:        "rejects an authenticated non-admin when auth is enabled",
			authEnabled: true,
			withCookie:  true,
			groups:      []string{"users"},
			wantStatus:  http.StatusForbidden,
		},
		{
			name:        "allows an authenticated admin when auth is enabled",
			authEnabled: true,
			withCookie:  true,
			groups:      []string{"speedtest-admins"},
			wantStatus:  http.StatusNoContent,
			wantCalled:  true,
		},
		{
			name:       "allows requests when auth is disabled",
			wantStatus: http.StatusNoContent,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			e := echo.New()
			e.Use(authapi.Middleware(authapi.MiddlewareOpts{
				Enabled:     tt.authEnabled,
				Service:     &notificationsSessionServiceStub{groups: tt.groups},
				AdminGroup:  "speedtest-admins",
				Routes:      authRoutes(),
				AdminRoutes: adminRoutes(),
			})...)
			e.GET("/v1/config/notifications", func(c *echo.Context) error {
				called = true
				return c.NoContent(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/v1/config/notifications", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "__Host-speedtest-session", Value: "session-id"})
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equalf(
				t,
				tt.wantStatus,
				rec.Code,
				"status = %d, want %d",
				rec.Code,
				tt.wantStatus,
			)
			assert.Equalf(
				t,
				tt.wantCalled,
				called,
				"handler called = %t, want %t",
				called,
				tt.wantCalled,
			)
		})
	}
}

type notificationsSessionServiceStub struct {
	groups []string
}

func (s *notificationsSessionServiceStub) Invoke(
	_ context.Context,
	_ string,
	_ string,
	_ string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	return notificationsSessionEnvelope{groups: s.groups}, nil
}

type notificationsSessionEnvelope struct {
	groups []string
}

func (e notificationsSessionEnvelope) Decode(into any) error {
	*into.(*authactor.AuthenticateResult) = authactor.AuthenticateResult{Groups: e.groups}
	return nil
}

func hasProtectedRoute(routes authapi.ProtectedRoutes, method, path string) bool {
	return slices.ContainsFunc(routes, func(route authapi.ProtectedRoute) bool {
		return route.Method == method && route.Path == path
	})
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		headers map[string]string
		status  int
	}{
		{
			name:   "configured UI origin",
			method: http.MethodGet,
			headers: map[string]string{
				echo.HeaderOrigin: "https://ui.example",
			},
			status: http.StatusOK,
		},
		{
			name:   "post preflight",
			method: http.MethodOptions,
			headers: map[string]string{
				echo.HeaderOrigin:                      "https://ui.example",
				echo.HeaderAccessControlRequestMethod:  http.MethodPost,
				echo.HeaderAccessControlRequestHeaders: echo.HeaderContentType,
			},
			status: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(corsMiddleware("https://ui.example"))
			e.GET("/v1/speedtest", func(c *echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/v1/speedtest", nil)
			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equalf(t, tt.status, rec.Code, "status = %d, want %d", rec.Code, tt.status)
			if origin := rec.Header().
				Get(echo.HeaderAccessControlAllowOrigin); origin != "https://ui.example" {
				assert.Failf(
					t,
					"assertion failed",
					"Access-Control-Allow-Origin = %q, want configured UI origin",
					origin,
				)
			}
			if credentials := rec.Header().
				Get(echo.HeaderAccessControlAllowCredentials); credentials != "true" {
				assert.Failf(
					t,
					"assertion failed",
					"Access-Control-Allow-Credentials = %q, want true",
					credentials,
				)
			}
		})
	}
}

func TestStaticMiddleware(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(
		filepath.Join(root, "index.html"),
		[]byte("<h1>Speedtest Tracker</h1>"),
		0o644,
	); err != nil {
		require.FailNow(t, "assertion failed", err)
	}

	e := echo.New()
	e.Filesystem = os.DirFS(root)
	e.Use(staticMiddleware("."))
	e.GET("/v1/speedtest", func(c *echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "serves the SPA for a client-side route",
			path:       "/results/01a16793-c56d-73f3-9530-352ea76047b2",
			wantStatus: http.StatusOK,
			wantBody:   "<h1>Speedtest Tracker</h1>",
		},
		{
			name:       "does not serve the SPA for an unknown API route",
			path:       "/v1/missing",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equalf(
				t,
				tt.wantStatus,
				rec.Code,
				"status = %d, want %d",
				rec.Code,
				tt.wantStatus,
			)
			if body := rec.Body.String(); tt.wantBody != "" && body != tt.wantBody {
				assert.Failf(t, "assertion failed", "body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
