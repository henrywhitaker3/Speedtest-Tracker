package speedtest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
)

func TestDeleteHandlerDeletesResult(t *testing.T) {
	id := uuid.MustOrdered()
	repo := &resultDeleteRepoStub{}

	_, err := (&deleteHandler{repo: repo}).Handler()(testContext(t), DeleteRequest{ID: id})
	require.NoError(t, err)
	assert.True(t, repo.called)
	assert.Equal(t, id, repo.id)
}

func TestDeleteHandlerWrapsRepositoryError(t *testing.T) {
	want := errors.New("delete result")
	repo := &resultDeleteRepoStub{err: want}

	_, err := (&deleteHandler{repo: repo}).Handler()(
		testContext(t),
		DeleteRequest{ID: uuid.MustOrdered()},
	)
	assert.ErrorIs(t, err, want)
}

func TestDeleteEndpointRequiresAuthenticatedAdministrator(t *testing.T) {
	tests := []struct {
		name       string
		withCookie bool
		groups     []string
		wantStatus int
		wantDelete bool
	}{
		{
			name:       "rejects unauthenticated requests",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects authenticated non-administrators",
			withCookie: true,
			groups:     []string{"users"},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "allows authenticated administrators",
			withCookie: true,
			groups:     []string{"administrators"},
			wantStatus: http.StatusOK,
			wantDelete: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &resultDeleteRepoStub{}
			h := windowhttp.New(windowhttp.HTTPOpts{})
			routes := authapi.ProtectedRoutes{{
				Method: http.MethodDelete,
				Path:   "/v1/speedtest/:id",
			}}
			for _, middleware := range authapi.Middleware(authapi.MiddlewareOpts{
				Enabled:     true,
				Service:     &authenticatedSessionInvokerStub{groups: tt.groups},
				AdminGroup:  "administrators",
				Routes:      routes,
				AdminRoutes: routes,
			}) {
				h.Use(middleware)
			}
			h.Register(&deleteHandler{repo: repo})

			req := httptest.NewRequest(
				http.MethodDelete,
				"/v1/speedtest/"+uuid.MustOrdered().String(),
				nil,
			)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "__Host-speedtest-session", Value: "session-id"})
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			assert.Equalf(
				t,
				tt.wantStatus,
				rec.Code,
				"status = %d, want %d",
				rec.Code,
				tt.wantStatus,
			)
			assert.Equalf(
				t,
				tt.wantDelete,
				repo.called,
				"Delete called = %t, want %t",
				repo.called,
				tt.wantDelete,
			)
		})
	}
}

type resultDeleteRepoStub struct {
	called bool
	id     uuid.UUID
	err    error
}

func (r *resultDeleteRepoStub) Delete(_ context.Context, id uuid.UUID) error {
	r.called = true
	r.id = id
	return r.err
}
