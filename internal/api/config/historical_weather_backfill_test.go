package config

import (
	"context"
	"errors"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/italypaleale/francis/actor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type backfillContextKey struct{}

func TestHistoricalWeatherBackfillHandler(t *testing.T) {
	dispatchErr := errors.New("dispatch failed")
	for _, tt := range []struct {
		name string
		err  error
	}{
		{name: "dispatches backfill"},
		{name: "propagates dispatch error", err: dispatchErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := &backfillDispatcherStub{err: tt.err}
			c := testContext(t)
			ctx := context.WithValue(c.Request().Context(), backfillContextKey{}, "request value")
			c.SetRequest(c.Request().WithContext(ctx))
			got, err := (&historicalWeatherBackfillHandler{dispatcher: stub}).Handler()(c, nil)
			assert.ErrorIs(t, err, tt.err)
			require.NotNil(t, stub.ctx)
			assert.Equal(t, "request value", stub.ctx.Value(backfillContextKey{}))
			assert.Equal(t, weather.Actor, stub.kind)
			assert.Equal(t, actor.SingletonActorID, stub.id)
			assert.Equal(t, weather.Backfill, stub.method)
			assert.Nil(t, stub.data)
			if tt.err == nil {
				require.NotNil(t, got)
				assert.Equal(
					t,
					HistoricalWeatherBackfillResponse{JobID: "job-1", Created: true},
					*got,
				)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

type backfillDispatcherStub struct {
	ctx              context.Context
	kind, id, method string
	data             any
	err              error
}

func (s *backfillDispatcherStub) Dispatch(
	ctx context.Context,
	kind, id, method string,
	data any,
	_ ...actor.JobOption,
) (string, bool, error) {
	s.ctx, s.kind, s.id, s.method, s.data = ctx, kind, id, method, data
	return "job-1", true, s.err
}
