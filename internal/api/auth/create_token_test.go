package auth

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

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

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	if repo.called {
		t.Error("repository Create was called for a blank name")
	}
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

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	if repo.called {
		t.Error("repository Create was called for a name shorter than five characters")
	}
}

func TestCreateTokenWithShortNameDoesNotPersistToken(t *testing.T) {
	db, err := database.Sqlite(filepath.Join(t.TempDir(), "auth.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := internalAuth.Register(db); err != nil {
		t.Fatal(err)
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

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}

	tokens, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(tokens) != 0 {
		t.Errorf("stored tokens = %d, want 0", len(tokens))
	}
}

func TestCreateTokenHandlerRequiresAuthentication(t *testing.T) {
	repo := &createRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Use(SessionMiddleware())
	h.Use(AuthMiddleware(AuthMiddlewareOpts{
		TokenRepo: repo,
		Routes: ProtectedRoutes{{
			Method: http.MethodPost,
			Path:   "/v1/auth/tokens",
		}},
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
