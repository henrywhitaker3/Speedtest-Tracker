package weather

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
