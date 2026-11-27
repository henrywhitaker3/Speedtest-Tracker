package weather

import (
	"context"
	"fmt"
)

type Provider interface {
	Get(ctx context.Context) (*WeatherData, error)
}

func BuildProvider(provider, apiKey string, lon, lat float64) (Provider, error) {
	switch provider {
	case "openweathermap":
		return newOpenWeatherMap(apiKey, lon, lat), nil
	default:
		return nil, fmt.Errorf("unknown provider %s", provider)
	}
}
