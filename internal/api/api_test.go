package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"testing"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
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
		{method: http.MethodGet, path: "/v1/speedtest/latest"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			if got := hasProtectedRoute(authRoutes(), tt.method, tt.path); got != tt.auth {
				t.Errorf("authenticated = %t, want %t", got, tt.auth)
			}
			if got := hasProtectedRoute(adminRoutes(), tt.method, tt.path); got != tt.adminOnly {
				t.Errorf("admin-only = %t, want %t", got, tt.adminOnly)
			}
		})
	}
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

			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d", rec.Code, tt.status)
			}
			if origin := rec.Header().Get(echo.HeaderAccessControlAllowOrigin); origin != "https://ui.example" {
				t.Errorf("Access-Control-Allow-Origin = %q, want configured UI origin", origin)
			}
			if credentials := rec.Header().Get(echo.HeaderAccessControlAllowCredentials); credentials != "true" {
				t.Errorf("Access-Control-Allow-Credentials = %q, want true", credentials)
			}
		})
	}
}

func TestStaticMiddleware(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "index.html"), []byte("<h1>Speedtest Tracker</h1>"), 0o644); err != nil {
		t.Fatal(err)
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

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if body := rec.Body.String(); tt.wantBody != "" && body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
