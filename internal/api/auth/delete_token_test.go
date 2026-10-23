package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
)

func TestDeleteTokenHandlerRequiresAuthentication(t *testing.T) {
	repo := &deleteTokenRepoStub{}
	h := windowhttp.New(windowhttp.HTTPOpts{})
	h.Use(SessionMiddleware())
	h.Use(AuthMiddleware(AuthMiddlewareOpts{
		TokenRepo: repo,
		Routes: ProtectedRoutes{{
			Method: http.MethodDelete,
			Path:   "/v1/auth/tokens/:id",
		}},
	}))
	h.Register(&deleteHandler{repo: repo})

	req := httptest.NewRequest(
		http.MethodDelete,
		"/v1/auth/tokens/"+uuid.MustOrdered().String(),
		nil,
	)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if repo.called {
		t.Error("repository Delete was called without authentication")
	}
}

type deleteTokenRepoStub struct {
	called bool
}

func (r *deleteTokenRepoStub) Validate(context.Context, string) error {
	return nil
}

func (r *deleteTokenRepoStub) Delete(context.Context, uuid.UUID) error {
	r.called = true
	return nil
}
