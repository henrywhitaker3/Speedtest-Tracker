package weather

import "context"

type Provider interface {
	Get(ctx context.Context, lon, lat float64) (*WeatherData, error)
}
