package weather

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
)

type HistoricalProvider interface {
	Get(context.Context, time.Time) (*WeatherData, error)
}

type ThrottledError struct {
	RetryAfter time.Duration
}

func (e *ThrottledError) Error() string {
	return fmt.Sprintf("weather provider throttled; retry after %s", e.RetryAfter)
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
		if len(toStore) >= 500 {
			if err := h.repo.Store(ctx, toStore); err != nil {
				return fmt.Errorf("store historical data: %w", err)
			}
			toStore = []Weather{}
		}
		val, err := h.getWeatherData(ctx, t)
		if err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("get historical weather data: %w", ctx.Err())
			}
			slog.ErrorContext(ctx, "get historical weather data", "error", err)
			continue
		}
		if val == nil {
			slog.ErrorContext(ctx, "historical weather provider returned no data", "result", id)
			continue
		}
		data := Weather{
			ResultID:    id,
			WeatherData: *val,
		}
		toStore = append(toStore, data)
		slog.DebugContext(ctx, "retrieved historical weather data", "data", data)
	}

	if len(toStore) == 0 {
		return nil
	}
	if err := h.repo.Store(ctx, toStore); err != nil {
		return fmt.Errorf("store historical data: %w", err)
	}

	return nil
}

func (h *historicalActor) getWeatherData(ctx context.Context, t time.Time) (*WeatherData, error) {
	for {
		val, err := h.provider.Get(ctx, t)
		if err == nil {
			return val, nil
		}
		var throttled *ThrottledError
		if !errors.As(err, &throttled) || throttled.RetryAfter <= 0 {
			return nil, err
		}
		timer := time.NewTimer(throttled.RetryAfter)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, fmt.Errorf("wait for weather provider: %w", ctx.Err())
		case <-timer.C:
		}
	}
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
