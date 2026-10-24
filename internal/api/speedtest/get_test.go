package speedtest

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/builtin/workflow"
)

func TestGetHandler(t *testing.T) {
	result := &results.Result{BaseModel: database.BaseModel{ID: uuid.MustOrdered()}, Complete: true}
	repoErr := errors.New("get failed")

	tests := []struct {
		name       string
		result     *results.Result
		repoErr    error
		status     workflow.InstanceStatus
		statusErr  error
		wantResult *GetResponse
		wantErr    error
	}{
		{
			name:       "returns result",
			result:     result,
			wantResult: &GetResponse{Result: Result{ID: result.ID, Complete: true}},
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
			status := &workflowCheckerStub{status: tt.status, err: tt.statusErr}
			handler := (&getHandler{repo: repo, status: status}).Handler()
			request := GetRequest{ID: uuid.MustOrdered()}

			got, err := handler(testContext(t), request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Handler() error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("Handler() result = %v, want %v", got, tt.wantResult)
			}
			if repo.id != request.ID {
				t.Errorf("repository ID = %q, want %q", repo.id, request.ID)
			}
		})
	}
}

type workflowCheckerStub struct {
	status workflow.InstanceStatus
	err    error
	id     string
}

func (w *workflowCheckerStub) GetStatus(_ context.Context, id string) (workflow.InstanceStatus, error) {
	w.id = id
	return w.status, w.err
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
