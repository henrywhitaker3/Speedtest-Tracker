package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenWeatherMapGet(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		want      *WeatherData
		wantError string
	}{
		{
			name:   "weather",
			status: 200,
			body:   `{"weather":[{"id":500}],"main":{"temp":12.5},"wind":{"speed":3.2},"rain":{"1h":1.1},"snow":{"1h":0.2}}`,
			want:   &WeatherData{Code: 61, Temperature: 12.5, WindSpeed: 3.2, Precipitation: 1.3},
		},
		{
			name:   "unknown code",
			status: 200,
			body:   `{"weather":[{"id":999}]}`,
			want:   &WeatherData{Code: -1},
		},
		{name: "missing conditions", status: 200, body: `{}`, wantError: "no weather conditions"},
		{name: "invalid response", status: 200, body: `{`, wantError: "parse response body"},
		{
			name:      "API error",
			status:    401,
			body:      `unauthorized`,
			wantError: "got status code 401: unauthorized",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "/data/2.5/weather", r.URL.Path)
					assert.Equal(t, "-0.12", r.URL.Query().Get("lon"))
					assert.Equal(t, "51.5", r.URL.Query().Get("lat"))
					assert.Equal(t, "test-key", r.URL.Query().Get("appid"))
					assert.Equal(t, "metric", r.URL.Query().Get("units"))
					w.WriteHeader(tt.status)
					_, _ = w.Write([]byte(tt.body))
				}),
			)
			defer server.Close()
			provider := newOpenWeatherMap("test-key")
			provider.baseURL = server.URL
			got, err := provider.Get(context.Background(), -0.12, 51.5)
			if tt.wantError != "" {
				require.ErrorContains(t, err, tt.wantError)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
