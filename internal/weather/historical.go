package weather

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
)

type HistoricalProvider interface {
	Get(context.Context, time.Time) (*WeatherData, error)
}

func BuildHistoricalProvider(
	provider, apiKey string,
	lon, lat float64,
) (HistoricalProvider, error) {
	switch provider {
	case "apifreaks":
		return newAPIFreaks(apiKey, lon, lat), nil
	default:
		return nil, fmt.Errorf("unknown provider %s", provider)
	}
}

type historicalActorState struct{}

type historicalWeatherRepo interface {
	WithoutWeather(context.Context) (map[uuid.UUID]time.Time, error)
	Store(context.Context, []Weather) error
}

type historicalActor struct {
	client   actor.Client[historicalActorState]
	repo     historicalWeatherRepo
	provider HistoricalProvider
}

const (
	Actor    = "historical-weather"
	Backfill = "backfill"
)

func (h *historicalActor) Job(ctx context.Context, method string, data actor.Envelope) error {
	switch method {
	case Backfill:
		return h.jobBackfill(ctx)
	default:
		return fmt.Errorf("unknown method %s: %w", method, actor.ErrJobPermanentFailure)
	}
}

func (h *historicalActor) jobBackfill(ctx context.Context) error {
	toBackfill, err := h.repo.WithoutWeather(ctx)
	if err != nil {
		return fmt.Errorf("get results with missing weather: %w", err)
	}

	toStore := []Weather{}
	for id, t := range toBackfill {
		val, err := h.provider.Get(ctx, t)
		if err != nil {
			return fmt.Errorf("get historical weather data: %w", err)
		}
		toStore = append(toStore, Weather{
			ResultID:    id,
			WeatherData: *val,
		})
	}

	if err := h.repo.Store(ctx, toStore); err != nil {
		return fmt.Errorf("store historical data: %w", err)
	}

	return nil
}

type store struct {
	mu   *sync.Mutex
	data map[time.Time]*WeatherData
}

func (s *store) get(t time.Time) (*WeatherData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, ok := s.data[normalisedTime(t)]
	return val, ok
}

func (s *store) put(t time.Time, d *WeatherData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[normalisedTime(t)] = d
}

func normalisedTime(t time.Time) time.Time {
	return t.UTC().Round(time.Hour)
}
