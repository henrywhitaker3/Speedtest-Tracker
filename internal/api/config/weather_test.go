package config

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCurrentWeatherHandler(t *testing.T) {
	data := &weather.WeatherData{Code: 61, Temperature: 12.5, Precipitation: 1.2, WindSpeed: 3.4}
	providerErr := errors.New("provider unavailable")
	tests := []struct {
		name     string
		provider *weatherProviderStub
		want     *weather.WeatherData
		wantErr  error
		disabled bool
	}{
		{name: "returns weather", provider: &weatherProviderStub{data: data}, want: data},
		{
			name:     "returns provider error",
			provider: &weatherProviderStub{err: providerErr},
			wantErr:  providerErr,
		},
		{name: "handles disabled weather", disabled: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var provider weather.Provider
			if tt.provider != nil {
				provider = tt.provider
			}
			c := testContext(t)
			got, err := (&getCurrentWeatherHandler{p: provider}).Handler()(c, nil)
			if tt.disabled {
				var httpErr *echo.HTTPError
				require.ErrorAs(t, err, &httpErr)
				assert.Equal(t, http.StatusServiceUnavailable, httpErr.Code)
			} else if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			if tt.provider != nil {
				assert.Equal(t, c.Request().Context(), tt.provider.ctx)
			}
		})
	}
}

type weatherProviderStub struct {
	data *weather.WeatherData
	err  error
	ctx  context.Context
}

func (p *weatherProviderStub) Get(ctx context.Context) (*weather.WeatherData, error) {
	p.ctx = ctx
	return p.data, p.err
}
