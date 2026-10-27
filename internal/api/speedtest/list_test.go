package speedtest

import (
	"context"
	"errors"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/database"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/assert"
)

func TestListHandler(t *testing.T) {
	page := &database.Pagination[results.Result]{
		Page:       2,
		Limit:      1,
		Total:      3,
		TotalPages: 3,
		Data: []results.Result{
			{ID: uuid.MustOrdered()},
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
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.page == nil {
				assert.Nilf(t, got, "Handler() result = %v, want nil", got)
			} else {
				want := &database.Pagination[Result]{
					Page:       tt.page.Page,
					Limit:      tt.page.Limit,
					Total:      tt.page.Total,
					TotalPages: tt.page.TotalPages,
				}
				for _, result := range tt.page.Data {
					want.Data = append(want.Data, Result{ID: result.ID})
				}
				assert.Equal(t, want, got)
			}
			assert.Equalf(
				t,
				tt.wantOpts,
				repo.opts,
				"repository options = %+v, want %+v",
				repo.opts,
				tt.wantOpts,
			)
		})
	}
}

type listRepoStub struct {
	page *database.Pagination[results.Result]
	err  error
	opts results.ListOpts
}

func (r *listRepoStub) List(
	_ context.Context,
	opts results.ListOpts,
) (*database.Pagination[results.Result], error) {
	r.opts = opts
	return r.page, r.err
}
