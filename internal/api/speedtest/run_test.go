package speedtest

import (
	"context"
	"errors"
	"testing"
	"time"

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

type actorInvokerStub struct {
	envelope  actor.Envelope
	err       error
	actorType string
	actorID   string
	method    string
	invoked   bool
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
