package speedtest

import (
	"context"
	"errors"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
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

type getRepoStub struct {
	result *results.Result
	err    error
	id     uuid.UUID
}

func (r *getRepoStub) Get(_ context.Context, id uuid.UUID) (*results.Result, error) {
	r.id = id
	return r.result, r.err
}
