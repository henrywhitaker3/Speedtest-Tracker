package weather

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBackfillSkipsFailedOrEmptyWeather(t *testing.T) {
	for _, tt := range []struct {
		name string
		data *WeatherData
		err  error
	}{
		{name: "provider error", err: errors.New("provider unavailable")},
		{name: "provider error with data", data: &WeatherData{}, err: errors.New("provider unavailable")},
		{name: "nil data"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			repo := &backfillRepoStub{id: uuid.MustOrdered()}
			err := (&historicalActor{repo: repo, provider: failingWeatherProvider{data: tt.data, err: tt.err}}).jobBackfill(context.Background())
			require.NoError(t, err)
			assert.Nil(t, repo.stored)
		})
	}
}

type failingWeatherProvider struct {
	data *WeatherData
	err  error
}

func (p failingWeatherProvider) Get(context.Context, time.Time) (*WeatherData, error) {
	return p.data, p.err
}

func TestBackfillRetriesThrottledProvider(t *testing.T) {
	id := uuid.MustOrdered()
	repo := &backfillRepoStub{id: id}
	provider := &throttledProviderStub{delay: time.Millisecond}
	err := (&historicalActor{repo: repo, provider: provider}).jobBackfill(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, provider.calls)
	require.Len(t, repo.stored, 1)
	assert.Equal(t, id, repo.stored[0].ResultID)
	assert.Equal(t, 3.0, repo.stored[0].Temperature)
}

func TestBackfillStopsWaitingOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := (&historicalActor{repo: &backfillRepoStub{id: uuid.MustOrdered()}, provider: &throttledProviderStub{delay: time.Hour}}).jobBackfill(ctx)
	assert.ErrorIs(t, err, context.Canceled)
}

type backfillRepoStub struct {
	id     uuid.UUID
	stored []Weather
}

func (r *backfillRepoStub) WithoutWeather(context.Context) (map[uuid.UUID]time.Time, error) {
	return map[uuid.UUID]time.Time{r.id: time.Now()}, nil
}
func (r *backfillRepoStub) Store(_ context.Context, data []Weather) error {
	r.stored = data
	return nil
}

type throttledProviderStub struct {
	calls int
	delay time.Duration
}

func (p *throttledProviderStub) Get(context.Context, time.Time) (*WeatherData, error) {
	p.calls++
	if p.calls == 1 {
		return nil, fmt.Errorf("wrapped: %w", &ThrottledError{RetryAfter: p.delay})
	}
	return &WeatherData{Temperature: 3}, nil
}

func TestNormalisedTime(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  time.Time
	}{
		{name: "rounds down", input: time.Date(2026, 11, 27, 10, 29, 59, 0, time.UTC), want: time.Date(2026, 11, 27, 10, 0, 0, 0, time.UTC)},
		{name: "rounds half hour up", input: time.Date(2026, 11, 27, 10, 30, 0, 0, time.UTC), want: time.Date(2026, 11, 27, 11, 0, 0, 0, time.UTC)},
		{name: "normalises timezone and date", input: time.Date(2026, 11, 28, 5, 50, 0, 0, time.FixedZone("offset", 5*60*60+30*60)), want: time.Date(2026, 11, 28, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, normalisedTime(tt.input))
		})
	}
}

func TestStoreUsesNormalisedTime(t *testing.T) {
	s := store{mu: &sync.Mutex{}, data: make(map[time.Time]*WeatherData)}
	data := &WeatherData{Temperature: 3}
	s.put(time.Date(2026, 11, 27, 10, 10, 0, 0, time.UTC), data)
	got, ok := s.get(time.Date(2026, 11, 27, 11, 20, 0, 0, time.FixedZone("offset", 60*60)))
	require.True(t, ok)
	assert.Same(t, data, got)
	_, ok = s.get(time.Date(2026, 11, 27, 10, 30, 0, 0, time.UTC))
	assert.False(t, ok)
}
