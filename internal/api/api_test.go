package api

import (
	"net/http"
	"net/http/httptest"
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
