package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	internalAuth "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
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

type createRepoStub struct {
	called bool
}

func (r *createRepoStub) Create(_ context.Context, _ string) (string, internalAuth.Token, error) {
	r.called = true
	return "", internalAuth.Token{}, nil
}
