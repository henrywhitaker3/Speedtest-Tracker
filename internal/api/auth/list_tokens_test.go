package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	internalAuth "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/stretchr/testify/assert"
)

func TestListTokensHandlerRequiresAuthentication(t *testing.T) {
	repo := &listTokensRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	routes := ProtectedRoutes{{
		Method: http.MethodGet,
		Path:   "/v1/auth/tokens",
	}}
	for _, middleware := range Middleware(MiddlewareOpts{
		Enabled:     true,
		TokenRepo:   repo,
		Routes:      routes,
		AdminRoutes: routes,
	}) {
		h.Use(middleware)
	}
	h.Register(&listTokensHandler{repo: repo})

	req := httptest.NewRequest(http.MethodGet, "/v1/auth/tokens", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equalf(
		t,
		http.StatusUnauthorized,
		rec.Code,
		"status = %d, want %d",
		rec.Code,
		http.StatusUnauthorized,
	)
	assert.False(t, repo.called)
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
