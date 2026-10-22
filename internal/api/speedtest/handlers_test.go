package speedtest

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/ookla"
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

func TestGetHandler(t *testing.T) {
	result := &results.Result{BaseModel: database.BaseModel{ID: uuid.MustOrdered()}, Complete: true}
	repoErr := errors.New("get failed")

	tests := []struct {
		name       string
		result     *results.Result
		repoErr    error
		wantResult *results.Result
		wantErr    error
	}{
		{
			name:       "returns result",
			result:     result,
			wantResult: result,
		},
		{
			name:    "wraps repository error",
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &getRepoStub{result: tt.result, err: tt.repoErr}
			handler := (&getHandler{repo: repo}).Handler()
			request := GetRequest{ID: uuid.MustOrdered()}

			got, err := handler(testContext(t), request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.wantResult {
				t.Errorf("Handler() result = %v, want %v", got, tt.wantResult)
			}
			if repo.id != request.ID {
				t.Errorf("repository ID = %q, want %q", repo.id, request.ID)
			}
		})
	}
}

func TestListHandler(t *testing.T) {
	page := &database.Pagination[results.Result]{
		Page:       2,
		Limit:      1,
		Total:      3,
		TotalPages: 3,
		Data: []results.Result{
			{BaseModel: database.BaseModel{ID: uuid.MustOrdered()}},
		},
	}
	repoErr := errors.New("list failed")

	tests := []struct {
		name     string
		opts     ListRequest
		page     *database.Pagination[results.Result]
		repoErr  error
		wantErr  error
		wantOpts results.ListOpts
	}{
		{
			name:     "returns paginated results",
			opts:     ListRequest{Page: 2, Limit: 1},
			page:     page,
			wantOpts: results.ListOpts{Page: 2, Limit: 1},
		},
		{
			name:     "wraps repository error",
			opts:     ListRequest{Page: 3, Limit: 10},
			repoErr:  repoErr,
			wantErr:  repoErr,
			wantOpts: results.ListOpts{Page: 3, Limit: 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &listRepoStub{page: tt.page, err: tt.repoErr}
			handler := (&listHandler{repo: repo}).Handler()

			got, err := handler(testContext(t), tt.opts)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if got != tt.page {
				t.Errorf("Handler() result = %v, want %v", got, tt.page)
			}
			if repo.opts != tt.wantOpts {
				t.Errorf("repository options = %+v, want %+v", repo.opts, tt.wantOpts)
			}
		})
	}
}

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

type getRepoStub struct {
	result *results.Result
	err    error
	id     uuid.UUID
}

func (r *getRepoStub) Get(_ context.Context, id uuid.UUID) (*results.Result, error) {
	r.id = id
	return r.result, r.err
}

type listRepoStub struct {
	page *database.Pagination[results.Result]
	err  error
	opts results.ListOpts
}

func (r *listRepoStub) List(_ context.Context, opts results.ListOpts) (*database.Pagination[results.Result], error) {
	r.opts = opts
	return r.page, r.err
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

func testContext(t *testing.T) *echo.Context {
	t.Helper()
	return echo.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
}
