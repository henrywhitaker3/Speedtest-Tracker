package speedtest

import (
	"context"
	"errors"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
)

func TestRunHandler(t *testing.T) {
	result := results.Result{BaseModel: database.BaseModel{ID: uuid.MustOrdered()}}
	invokeErr := errors.New("invoke failed")
	decodeErr := errors.New("decode failed")

	tests := []struct {
		name      string
		invokeErr error
		decodeErr error
		wantID    uuid.UUID
		wantErr   error
	}{
		{
			name:   "returns actor result ID",
			wantID: result.ID,
		},
		{
			name:      "wraps actor error",
			invokeErr: invokeErr,
			wantErr:   invokeErr,
		},
		{
			name:      "wraps decode error",
			decodeErr: decodeErr,
			wantErr:   decodeErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &actorInvokerStub{
				envelope: envelopeStub{result: result, err: tt.decodeErr},
				err:      tt.invokeErr,
			}
			handler := (&runHandler{service: service}).Handler()

			got, err := handler(testContext(t), nil)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (got == nil || got.ID != tt.wantID) {
				t.Errorf("Handler() result = %v, want ID %q", got, tt.wantID)
			}
			if service.actorType != ookla.Actor {
				t.Errorf("actor type = %q, want %q", service.actorType, ookla.Actor)
			}
			if service.method != ookla.Onboard {
				t.Errorf("actor method = %q, want %q", service.method, ookla.Onboard)
			}
			if service.actorID == "" {
				t.Error("actor ID is empty")
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
}

func (s *actorInvokerStub) Invoke(
	_ context.Context,
	actorType, actorID, method string,
	_ any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	s.actorType = actorType
	s.actorID = actorID
	s.method = method
	return s.envelope, s.err
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
