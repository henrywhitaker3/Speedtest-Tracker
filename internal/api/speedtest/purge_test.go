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
	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
)

func TestPurgeHandlerPurgesErroredResults(t *testing.T) {
	repo := &purgeRepoStub{}

	_, err := (&purgeHandler{repo: repo}).Handler()(testContext(t), nil)
	require.NoError(t, err)
	if !repo.called {
		assert.Fail(t, "assertion failed", "PurgeErrors was not called")
	}
}

func TestPurgeEndpointRequiresAuthenticatedAdministrator(t *testing.T) {
	tests := []struct {
		name       string
		withCookie bool
		groups     []string
		wantStatus int
		wantPurge  bool
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
			wantPurge:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &purgeRepoStub{}
			service := &authenticatedSessionInvokerStub{groups: tt.groups}
			h := windowhttp.New(windowhttp.HTTPOpts{})
			routes := authapi.ProtectedRoutes{{
				Method: http.MethodDelete,
				Path:   "/v1/speedtest/purge",
			}}
			for _, middleware := range authapi.Middleware(authapi.MiddlewareOpts{
				Enabled:     true,
				Service:     service,
				AdminGroup:  "administrators",
				Routes:      routes,
				AdminRoutes: routes,
			}) {
				h.Use(middleware)
			}
			h.Register(&purgeHandler{repo: repo})

			req := httptest.NewRequest(http.MethodDelete, "/v1/speedtest/purge", nil)
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
				tt.wantPurge,
				repo.called,
				"PurgeErrors called = %t, want %t",
				repo.called,
				tt.wantPurge,
			)
		})
	}
}

func TestPurgeHandlerWrapsRepositoryError(t *testing.T) {
	want := errors.New("delete results")
	repo := &purgeRepoStub{err: want}

	_, err := (&purgeHandler{repo: repo}).Handler()(testContext(t), nil)
	assert.ErrorIs(t, err, want)
}

type purgeRepoStub struct {
	called bool
	err    error
}

func (r *purgeRepoStub) PurgeErrors(context.Context) error {
	r.called = true
	return r.err
}

type authenticatedSessionInvokerStub struct {
	groups []string
}

func (s *authenticatedSessionInvokerStub) Invoke(
	_ context.Context,
	actorType, _ string,
	method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	if actorType != authactor.Actor || method != authactor.Session {
		return nil, errors.New("unexpected actor invocation")
	}
	return authenticatedSessionEnvelope{groups: s.groups}, nil
}

type authenticatedSessionEnvelope struct {
	groups []string
}

func (e authenticatedSessionEnvelope) Decode(into any) error {
	*into.(*authactor.AuthenticateResult) = authactor.AuthenticateResult{Groups: e.groups}
	return nil
}
