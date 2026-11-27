package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIFreaksHistoricalWeather(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		assert.Equal(t, "/v1.0/weather/historical", r.URL.Path)
		assert.Equal(t, "-2.25", r.URL.Query().Get("long"))
		assert.Equal(t, "53.45", r.URL.Query().Get("lat"))
		assert.Equal(t, "2024-10-15", r.URL.Query().Get("date"))
		assert.Equal(t, "hourly", r.URL.Query().Get("precision"))
		assert.Equal(t, "UTC", r.URL.Query().Get("timezone"))
		assert.Equal(t, "test-key", r.URL.Query().Get("apiKey"))
		_, _ = w.Write(
			[]byte(
				`{"historical":{"hourly":[{"timestamp":"2024-10-15T10:00","weather_code":61,"temperature_2m":3.5,"precipitation":1.2,"wind_speed_10m":18},{"timestamp":"2024-10-15T11:00","weather_code":3,"temperature_2m":4,"precipitation":0,"wind_speed_10m":7.2}]}}`,
			),
		)
	}))
	defer server.Close()
	provider := newAPIFreaks("test-key", -2.25, 53.45)
	provider.baseURL = server.URL

	got, err := provider.Get(context.Background(), time.Date(2024, 10, 15, 10, 15, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, &WeatherData{Code: 61, Temperature: 3.5, Precipitation: 1.2, WindSpeed: 5}, got)
	got, err = provider.Get(context.Background(), time.Date(2024, 10, 15, 11, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, &WeatherData{Code: 3, Temperature: 4, WindSpeed: 2}, got)
	assert.Equal(t, 1, calls)
}

func TestAPIFreaksLoadDailyDataErrors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{name: "API error", status: 401, body: "unauthorized", want: "status code 401"},
		{name: "invalid JSON", status: 200, body: "{", want: "parse response"},
		{name: "missing hourly data", status: 200, body: `{}`, want: "no hourly weather data"},
		{name: "invalid timestamp", status: 200, body: `{"historical":{"hourly":[{"timestamp":"2024-10-15T10:00"},{"timestamp":"invalid"}]}}`, want: "parse hourly timestamp"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(tt.status)
					_, _ = w.Write([]byte(tt.body))
				}),
			)
			defer server.Close()
			provider := newAPIFreaks("test-key", 0, 0)
			provider.baseURL = server.URL
			err := provider.loadDailyData(
				context.Background(),
				time.Date(2024, 10, 15, 10, 0, 0, 0, time.UTC),
			)
			require.ErrorContains(t, err, tt.want)
			assert.Empty(t, provider.store.data)
		})
	}
}

func TestAPIFreaksThrottle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(
			[]byte(
				`{"message":"Minutely API request limit exceeded. Please try again in one minute."}`,
			),
		)
	}))
	defer server.Close()
	provider := newAPIFreaks("test-key", 0, 0)
	provider.baseURL = server.URL
	_, err := provider.Get(context.Background(), time.Now())
	var throttled *ThrottledError
	require.ErrorAs(t, err, &throttled)
	assert.Equal(t, time.Minute, throttled.RetryAfter)
}
