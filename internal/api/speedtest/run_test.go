package speedtest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	authapi "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/api/auth"
	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/labstack/echo/v5"
)

func TestRunHandler(t *testing.T) {
	result := &results.Result{ID: uuid.MustOrdered()}
	createErr := errors.New("create result")
	startErr := errors.New("start workflow")
	limiterErr := errors.New("limiter failed")

	tests := []struct {
		name           string
		limiterAllowed bool
		limiterErr     error
		createErr      error
		startErr       error
		wantID         uuid.UUID
		wantErr        error
		wantCreate     bool
		wantStart      bool
		wantDelete     bool
	}{
		{name: "creates a result and starts its workflow", limiterAllowed: true, wantID: result.ID, wantCreate: true, wantStart: true},
		{name: "wraps result creation errors", limiterAllowed: true, createErr: createErr, wantErr: createErr, wantCreate: true},
		{name: "wraps workflow start errors and deletes the result", limiterAllowed: true, startErr: startErr, wantErr: startErr, wantCreate: true, wantStart: true, wantDelete: true},
		{name: "rejects rate limited request", wantErr: echo.ErrTooManyRequests},
		{name: "rejects limiter error", limiterErr: limiterErr, wantErr: echo.ErrTooManyRequests},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &resultCreatorStub{result: result, err: tt.createErr}
			workflow := &workflowRunnerStub{err: tt.startErr}
			limiter := &rateLimiterStub{allowed: tt.limiterAllowed, err: tt.limiterErr}
			handler := (&runHandler{repo: repo, workflow: workflow, limiter: limiter}).Handler()

			got, err := handler(testContext(t), nil)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (got == nil || got.ID != tt.wantID) {
				t.Errorf("Handler() result = %v, want ID %q", got, tt.wantID)
			}
			if repo.called != tt.wantCreate {
				t.Errorf("Create called = %t, want %t", repo.called, tt.wantCreate)
			}
			if workflow.called != tt.wantStart {
				t.Errorf("workflow Start called = %t, want %t", workflow.called, tt.wantStart)
			}
			if repo.deleteCalled != tt.wantDelete {
				t.Errorf("Delete called = %t, want %t", repo.deleteCalled, tt.wantDelete)
			}
			if tt.wantDelete && repo.deletedID != result.ID {
				t.Errorf("deleted result ID = %q, want %q", repo.deletedID, result.ID)
			}
			if tt.wantStart {
				input, ok := workflow.input.(ookla.RunWorkflow)
				if !ok || input.Result.ID != result.ID {
					t.Errorf("workflow input = %#v, want result %q", workflow.input, result.ID)
				}
			}
			if limiter.key != "test" {
				t.Errorf("limiter key = %q, want %q", limiter.key, "test")
			}
		})
	}
}

func TestRunEndpointAuthentication(t *testing.T) {
	result := &results.Result{ID: uuid.MustOrdered()}
	tests := []struct {
		name            string
		withCookie      bool
		authorization   string
		sessionErr      error
		tokenErr        error
		wantStatus      int
		wantRunCall     bool
		wantSessionCall bool
		wantTokenCall   bool
	}{
		{name: "rejects a request without a browser session cookie", wantStatus: http.StatusUnauthorized},
		{name: "rejects a browser session without authenticated state", withCookie: true, sessionErr: errors.New("no authenticated session"), wantStatus: http.StatusUnauthorized, wantSessionCall: true},
		{name: "allows an authenticated session", withCookie: true, wantStatus: http.StatusOK, wantRunCall: true, wantSessionCall: true},
		{name: "rejects an invalid API token", authorization: "Bearer invalid-token", tokenErr: errors.New("token not found"), wantStatus: http.StatusUnauthorized, wantTokenCall: true},
		{name: "allows a valid API token without a browser session", authorization: "Bearer valid-token", wantStatus: http.StatusOK, wantRunCall: true, wantTokenCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &runEndpointInvokerStub{sessionErr: tt.sessionErr}
			tokenRepo := &tokenRepoStub{err: tt.tokenErr}
			runner := &workflowRunnerStub{}
			run := (&runHandler{
				repo:     &resultCreatorStub{result: result},
				workflow: runner,
				limiter:  &rateLimiterStub{allowed: true},
			}).Handler()

			e := echo.New()
			e.Use(authapi.Middleware(authapi.MiddlewareOpts{
				Enabled:   true,
				Service:   service,
				TokenRepo: tokenRepo,
				Routes: authapi.ProtectedRoutes{{
					Method: http.MethodPost,
					Path:   "/v1/speedtest/run",
				}},
			})...)
			e.POST("/v1/speedtest/run", func(c *echo.Context) error {
				response, err := run(c, nil)
				if err != nil {
					return err
				}
				return c.JSON(http.StatusOK, response)
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/speedtest/run", nil)
			if tt.withCookie {
				req.AddCookie(&http.Cookie{Name: "__Host-speedtest-session", Value: "session-id"})
			}
			if tt.authorization != "" {
				req.Header.Set(echo.HeaderAuthorization, tt.authorization)
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if runner.called != tt.wantRunCall {
				t.Errorf("workflow Start called = %t, want %t", runner.called, tt.wantRunCall)
			}
			if service.sessionCalled != tt.wantSessionCall {
				t.Errorf("session actor called = %t, want %t", service.sessionCalled, tt.wantSessionCall)
			}
			if tokenRepo.called != tt.wantTokenCall {
				t.Errorf("token repository called = %t, want %t", tokenRepo.called, tt.wantTokenCall)
			}
			if tt.wantTokenCall && tokenRepo.token != strings.TrimPrefix(tt.authorization, "Bearer ") {
				t.Errorf("token = %q, want bearer token from request", tokenRepo.token)
			}
		})
	}
}

type resultCreatorStub struct {
	result *results.Result
	err    error
	called bool

	deleteCalled bool
	deletedID    uuid.UUID
}

func (r *resultCreatorStub) Create(context.Context) (*results.Result, error) {
	r.called = true
	return r.result, r.err
}

func (r *resultCreatorStub) Delete(_ context.Context, id uuid.UUID) error {
	r.deleteCalled = true
	r.deletedID = id
	return nil
}

type workflowRunnerStub struct {
	err    error
	input  any
	called bool
}

func (r *workflowRunnerStub) Start(_ context.Context, input any, _ ...workflow.StartOption) (string, bool, error) {
	r.called = true
	r.input = input
	return "workflow-id", true, r.err
}

type runEndpointInvokerStub struct {
	sessionErr    error
	sessionCalled bool
}

func (s *runEndpointInvokerStub) Invoke(
	_ context.Context,
	actorType, _ string,
	method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	if actorType != authactor.Actor || method != authactor.Session {
		return nil, errors.New("unexpected actor invocation")
	}
	s.sessionCalled = true
	return runSessionEnvelope{}, s.sessionErr
}

type runSessionEnvelope struct{}

func (runSessionEnvelope) Decode(into any) error {
	*into.(*authactor.AuthenticateResult) = authactor.AuthenticateResult{}
	return nil
}

type tokenRepoStub struct {
	err    error
	token  string
	called bool
}

func (s *tokenRepoStub) Validate(_ context.Context, token string) error {
	s.called = true
	s.token = token
	return s.err
}

type rateLimiterStub struct {
	allowed bool
	err     error
	key     string
}

func (s *rateLimiterStub) Allow(_ context.Context, key string) (bool, time.Duration, error) {
	s.key = key
	return s.allowed, 0, s.err
}
