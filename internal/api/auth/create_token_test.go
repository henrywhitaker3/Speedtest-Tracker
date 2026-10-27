package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	internalAuth "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
)

func TestCreateTokenRejectsBlankNameBeforeCreatingToken(t *testing.T) {
	repo := &createRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Register(&createHandler{repo: repo})

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/auth/tokens",
		bytes.NewBufferString(`{"name":""}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equalf(
		t,
		http.StatusUnprocessableEntity,
		rec.Code,
		"status = %d, want %d",
		rec.Code,
		http.StatusUnprocessableEntity,
	)
	assert.False(t, repo.called)
}

func TestCreateTokenRejectsShortNameBeforeCreatingToken(t *testing.T) {
	repo := &createRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Register(&createHandler{repo: repo})

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/auth/tokens",
		bytes.NewBufferString(`{"name":"test"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	assert.Equalf(
		t,
		http.StatusUnprocessableEntity,
		rec.Code,
		"status = %d, want %d",
		rec.Code,
		http.StatusUnprocessableEntity,
	)
	assert.False(t, repo.called)
}

func TestCreateTokenWithShortNameDoesNotPersistToken(t *testing.T) {
	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	require.NoError(t, err)
	if err := internalAuth.Register(db); err != nil {
		require.FailNow(t, "assertion failed", err)
	}

	repo := internalAuth.NewRepo(db)
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Register(&createHandler{repo: repo})

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/auth/tokens",
		bytes.NewBufferString(`{"name":"test"}`),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	require.Equalf(
		t,
		http.StatusUnprocessableEntity,
		rec.Code,
		"status = %d, want %d",
		rec.Code,
		http.StatusUnprocessableEntity,
	)

	tokens, err := repo.List(context.Background())
	require.Equalf(t, nil, err, "List() error = %v", err)
	require.Len(t, tokens, 0)
}

func TestCreateTokenHandlerRequiresAuthentication(t *testing.T) {
	repo := &createRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	routes := ProtectedRoutes{{
		Method: http.MethodPost,
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
	h.Register(&createHandler{repo: repo})

	req := httptest.NewRequest(http.MethodPost, "/v1/auth/tokens", nil)
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
