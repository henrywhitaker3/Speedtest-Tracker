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
	"github.com/labstack/echo/v5"
)

func TestRunHandler(t *testing.T) {
	result := results.Result{ID: uuid.MustOrdered()}
	invokeErr := errors.New("invoke failed")
	decodeErr := errors.New("decode failed")
	limiterErr := errors.New("limiter failed")

	tests := []struct {
		name           string
		limiterAllowed bool
		limiterErr     error
		invokeErr      error
		decodeErr      error
		wantID         uuid.UUID
		wantErr        error
		wantInvoked    bool
	}{
		{
			name:           "returns actor result ID",
			limiterAllowed: true,
			wantID:         result.ID,
			wantInvoked:    true,
		},
		{
			name:           "wraps actor error",
			limiterAllowed: true,
			invokeErr:      invokeErr,
			wantErr:        invokeErr,
			wantInvoked:    true,
		},
		{
			name:           "wraps decode error",
			limiterAllowed: true,
			decodeErr:      decodeErr,
			wantErr:        decodeErr,
			wantInvoked:    true,
		},
		{
			name:    "rejects rate limited request",
			wantErr: echo.ErrTooManyRequests,
		},
		{
			name:       "rejects limiter error",
			limiterErr: limiterErr,
			wantErr:    echo.ErrTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &actorInvokerStub{
				envelope: envelopeStub{result: result, err: tt.decodeErr},
				err:      tt.invokeErr,
			}
			limiter := &rateLimiterStub{allowed: tt.limiterAllowed, err: tt.limiterErr}
			handler := (&runHandler{service: service, limiter: limiter}).Handler()

			got, err := handler(testContext(t), nil)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (got == nil || got.ID != tt.wantID) {
				t.Errorf("Handler() result = %v, want ID %q", got, tt.wantID)
			}
			if tt.wantInvoked && service.actorType != ookla.Actor {
				t.Errorf("actor type = %q, want %q", service.actorType, ookla.Actor)
			}
			if tt.wantInvoked && service.method != ookla.Onboard {
				t.Errorf("actor method = %q, want %q", service.method, ookla.Onboard)
			}
			if service.actorID == "" {
				if tt.wantInvoked {
					t.Error("actor ID is empty")
				}
			}
			if service.invoked != tt.wantInvoked {
				t.Errorf("actor invoked = %t, want %t", service.invoked, tt.wantInvoked)
			}
			if limiter.key != "test" {
				t.Errorf("limiter key = %q, want %q", limiter.key, "test")
			}
		})
	}
}

func TestRunEndpointAuthentication(t *testing.T) {
	result := results.Result{ID: uuid.MustOrdered()}

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
		{
			name:       "rejects a request without a browser session cookie",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:            "rejects a browser session without authenticated state",
			withCookie:      true,
			sessionErr:      errors.New("no authenticated session"),
			wantStatus:      http.StatusUnauthorized,
			wantSessionCall: true,
		},
		{
			name:            "allows an authenticated session",
			withCookie:      true,
			wantStatus:      http.StatusOK,
			wantRunCall:     true,
			wantSessionCall: true,
		},
		{
			name:          "rejects an invalid API token",
			authorization: "Bearer invalid-token",
			tokenErr:      errors.New("token not found"),
			wantStatus:    http.StatusUnauthorized,
			wantTokenCall: true,
		},
		{
			name:          "allows a valid API token without a browser session",
			authorization: "Bearer valid-token",
			wantStatus:    http.StatusOK,
			wantRunCall:   true,
			wantTokenCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &runEndpointInvokerStub{
				sessionErr: tt.sessionErr,
				envelope:   envelopeStub{result: result},
			}
			tokenRepo := &tokenRepoStub{err: tt.tokenErr}
			run := (&runHandler{
				service: service,
				limiter: &rateLimiterStub{allowed: true},
			}).Handler()

			e := echo.New()
			e.Use(authapi.SessionMiddleware())
			e.Use(authapi.AuthMiddleware(authapi.AuthMiddlewareOpts{
				Service:   service,
				TokenRepo: tokenRepo,
				Routes: authapi.ProtectedRoutes{{
					Method: http.MethodPost,
					Path:   "/v1/speedtest/run",
				}},
			}))
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
			if service.runCalled != tt.wantRunCall {
				t.Errorf("run actor called = %t, want %t", service.runCalled, tt.wantRunCall)
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

type actorInvokerStub struct {
	envelope  actor.Envelope
	err       error
	actorType string
	actorID   string
	method    string
	invoked   bool
}

type runEndpointInvokerStub struct {
	sessionErr    error
	envelope      actor.Envelope
	runCalled     bool
	sessionCalled bool
}

func (s *runEndpointInvokerStub) Invoke(
	_ context.Context,
	actorType, _ string, method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	if actorType == authactor.Actor && method == authactor.Session {
		s.sessionCalled = true
		return nil, s.sessionErr
	}
	if actorType == ookla.Actor && method == ookla.Onboard {
		s.runCalled = true
		return s.envelope, nil
	}
	return nil, errors.New("unexpected actor invocation")
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

func (s *actorInvokerStub) Invoke(
	_ context.Context,
	actorType, actorID, method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	s.invoked = true
	s.actorType = actorType
	s.actorID = actorID
	s.method = method
	return s.envelope, s.err
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

type envelopeStub struct {
	result results.Result
	err    error
}

func (e envelopeStub) Decode(into any) error {
	if e.err != nil {
		return e.err
	}
	*into.(*results.Result) = e.result
	return nil
}
