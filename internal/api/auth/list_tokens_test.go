package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	internalAuth "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
)

func TestListTokensHandlerRequiresAuthentication(t *testing.T) {
	repo := &listTokensRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Use(SessionMiddleware())
	h.Use(AuthMiddleware(AuthMiddlewareOpts{
		TokenRepo: repo,
		Routes:    []string{"/v1/auth/tokens"},
	}))
	h.Register(&listTokensHandler{repo: repo})

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/tokens", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if repo.called {
		t.Error("repository List was called without authentication")
	}
}

type listTokensRepoStub struct {
	called bool
}

func (r *listTokensRepoStub) Validate(context.Context, string) error {
	return nil
}

func (r *listTokensRepoStub) List(context.Context) ([]internalAuth.Token, error) {
	r.called = true
	return nil, nil
}
