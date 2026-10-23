package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalAuth "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/henrywhitaker3/windowframe/v2/http/validation"
	"github.com/labstack/echo/v5"
)

func TestCreateTokenHandlerRejectsBlankName(t *testing.T) {
	repo := &createRepoStub{}
	handler := (&createHandler{repo: repo}).Handler()
	ctx := echo.NewContext(
		httptest.NewRequest(http.MethodPost, "/v1/auth/tokens", nil),
		httptest.NewRecorder(),
	)

	_, err := handler(ctx, CreateTokenRequest{Name: " \t "})
	var validationErr *validation.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("Handler() error = %v, want validation error", err)
	}
	if repo.called {
		t.Error("repository Create was called for a blank name")
	}
}

func TestCreateTokenHandlerRequiresAuthentication(t *testing.T) {
	repo := &createRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Use(SessionMiddleware())
	h.Use(AuthMiddleware(AuthMiddlewareOpts{
		TokenRepo: repo,
		Routes:    []string{"/v1/auth/tokens"},
	}))
	h.Register(&createHandler{repo: repo})

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/tokens", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if repo.called {
		t.Error("repository Create was called without authentication")
	}
}

type createRepoStub struct {
	called bool
}

func (r *createRepoStub) Create(_ context.Context, _ string) (string, internalAuth.Token, error) {
	r.called = true
	return "", internalAuth.Token{}, nil
}

func (r *createRepoStub) Validate(context.Context, string) error {
	return nil
}
