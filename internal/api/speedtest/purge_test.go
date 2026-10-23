package speedtest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	windowhttp "github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
)

func TestPurgeHandlerPurgesErroredResults(t *testing.T) {
	repo := &purgeRepoStub{}

	_, err := (&purgeHandler{repo: repo}).Handler()(testContext(t), nil)
	if err != nil {
		t.Fatal(err)
	}
	if !repo.called {
		t.Error("PurgeErrors was not called")
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
			service := &purgeSessionInvokerStub{groups: tt.groups}
			h := windowhttp.New(windowhttp.HTTPOpts{})
			h.Use(authapi.SessionMiddleware())
			h.Use(authapi.AuthMiddleware(authapi.AuthMiddlewareOpts{
				Service:    service,
				AdminGroup: "administrators",
				Routes:     []string{"/v1/speedtest/purge"},
			}))
			h.Use(authapi.Admin(authapi.AdminOpts{
				Routes: []string{"/v1/speedtest/purge"},
			}))
			h.Register(&purgeHandler{repo: repo})

			req := httptest.NewRequest(http.MethodDelete, "/v1/speedtest/purge", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "__Host-speedtest-session", Value: "session-id"})
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if repo.called != tt.wantPurge {
				t.Errorf("PurgeErrors called = %t, want %t", repo.called, tt.wantPurge)
			}
		})
	}
}

func TestPurgeHandlerWrapsRepositoryError(t *testing.T) {
	want := errors.New("delete results")
	repo := &purgeRepoStub{err: want}

	_, err := (&purgeHandler{repo: repo}).Handler()(testContext(t), nil)
	if !errors.Is(err, want) {
		t.Errorf("Handler() error = %v, want %v", err, want)
	}
}

type purgeRepoStub struct {
	called bool
	err    error
}

func (r *purgeRepoStub) PurgeErrors(context.Context) error {
	r.called = true
	return r.err
}

type purgeSessionInvokerStub struct {
	groups []string
}

func (s *purgeSessionInvokerStub) Invoke(
	_ context.Context,
	actorType, _ string,
	method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	if actorType != authactor.Actor || method != authactor.Session {
		return nil, errors.New("unexpected actor invocation")
	}
	return purgeSessionEnvelope{groups: s.groups}, nil
}

type purgeSessionEnvelope struct {
	groups []string
}

func (e purgeSessionEnvelope) Decode(into any) error {
	*into.(*authactor.AuthenticateResult) = authactor.AuthenticateResult{Groups: e.groups}
	return nil
}
