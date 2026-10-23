package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v5"
)

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		headers map[string]string
		status  int
	}{
		{
			name:   "cross-origin request",
			method: http.MethodGet,
			headers: map[string]string{
				echo.HeaderOrigin: "https://any-deployment.example",
			},
			status: http.StatusOK,
		},
		{
			name:   "post preflight",
			method: http.MethodOptions,
			headers: map[string]string{
				echo.HeaderOrigin:                      "https://any-deployment.example",
				echo.HeaderAccessControlRequestMethod:  http.MethodPost,
				echo.HeaderAccessControlRequestHeaders: echo.HeaderContentType,
			},
			status: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(corsMiddleware())
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
			if origin := rec.Header().Get(echo.HeaderAccessControlAllowOrigin); origin != "*" {
				t.Errorf("Access-Control-Allow-Origin = %q, want *", origin)
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
