package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/labstack/echo/v5"
)

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		withCookie bool
		groups     []string
		wantStatus int
		wantCalled bool
	}{
		{
			name:       "rejects an unauthenticated request",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects an authenticated non-admin",
			withCookie: true,
			groups:     []string{"users"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "allows an authenticated admin",
			withCookie: true,
			groups:     []string{"speedtest-admins", "users"},
			wantStatus: http.StatusNoContent,
			wantCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &actorInvokerStub{
				envelope: sessionEnvelope{value: authactor.AuthenticateResult{
					Subject: "issuer-user-id",
					Groups:  tt.groups,
				}},
			}
			called := false
			e := echo.New()
			routes := ProtectedRoutes{{
				Method: http.MethodGet,
				Path:   "/admin-only",
			}}
			e.Use(Middleware(MiddlewareOpts{
				Enabled:     true,
				Service:     service,
				AdminGroup:  "speedtest-admins",
				Routes:      routes,
				AdminRoutes: routes,
			})...)
			e.GET("/admin-only", func(c *echo.Context) error {
				called = true
				return c.NoContent(http.StatusNoContent)
			})

			req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-id"})
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if called != tt.wantCalled {
				t.Errorf("handler called = %t, want %t", called, tt.wantCalled)
			}
		})
	}
}

func TestAdminMiddlewareRejectsAPIToken(t *testing.T) {
	tokenRepo := &adminTokenRepoStub{}
	called := false
	e := echo.New()
	routes := ProtectedRoutes{{
		Method: http.MethodGet,
		Path:   "/admin-only",
	}}
	e.Use(Middleware(MiddlewareOpts{
		Enabled:     true,
		TokenRepo:   tokenRepo,
		Routes:      routes,
		AdminRoutes: routes,
	})...)
	e.GET("/admin-only", func(c *echo.Context) error {
		called = true
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer valid-api-token")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if !tokenRepo.called || tokenRepo.token != "valid-api-token" {
		t.Errorf(
			"token validation = called %t with %q, want valid API token",
			tokenRepo.called,
			tokenRepo.token,
		)
	}
	if called {
		t.Error("admin handler was called for an API token")
	}
}

func TestAuthMiddlewareOptionallyAuthenticatesPublicRoute(t *testing.T) {
	tests := []struct {
		name       string
		withCookie bool
		groups     []string
		wantAdmin  bool
	}{
		{
			name: "allows an anonymous request",
		},
		{
			name:       "identifies an administrator",
			withCookie: true,
			groups:     []string{"speedtest-admins"},
			wantAdmin:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(Middleware(MiddlewareOpts{
				Enabled: true,
				Service: &actorInvokerStub{
					envelope: sessionEnvelope{
						value: authactor.AuthenticateResult{Groups: tt.groups},
					},
				},
				AdminGroup: "speedtest-admins",
			})...)
			e.GET("/results", func(c *echo.Context) error {
				if got := IsAdmin(c); got != tt.wantAdmin {
					t.Errorf("IsAdmin() = %t, want %t", got, tt.wantAdmin)
				}
				return c.NoContent(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/results", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-id"})
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
			}
		})
	}
}

type adminTokenRepoStub struct {
	called bool
	token  string
}

func (r *adminTokenRepoStub) Validate(_ context.Context, token string) error {
	r.called = true
	r.token = token
	return nil
}
