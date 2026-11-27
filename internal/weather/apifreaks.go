package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type apiFreaks struct {
	apiKey   string
	baseURL  string
	lon, lat float64
	store    *store
}

const (
	apiFreaksBaseURL = "https://api.apifreaks.com"
)

func newAPIFreaks(apiKey string, lon, lat float64) apiFreaks {
	return apiFreaks{
		apiKey:  apiKey,
		baseURL: apiFreaksBaseURL,
		lon:     lon,
		lat:     lat,
		store: &store{
			mu:   &sync.Mutex{},
			data: map[time.Time]*WeatherData{},
		},
	}
}

func (a apiFreaks) Get(ctx context.Context, t time.Time) (*WeatherData, error) {
	if val, ok := a.store.get(t); ok {
		return val, nil
	}
	if err := a.loadDailyData(ctx, normalisedTime(t)); err != nil {
		return nil, fmt.Errorf("load historical weather: %w", err)
	}
	if val, ok := a.store.get(t); ok {
		return val, nil
	}
	return nil, fmt.Errorf("historical weather unavailable for %s", normalisedTime(t).Format(time.RFC3339))
}

func (a apiFreaks) loadDailyData(ctx context.Context, t time.Time) error {
	endpoint, err := url.Parse(a.baseURL + "/v1.0/weather/historical")
	if err != nil {
		return fmt.Errorf("build url: %w", err)
	}
	query := endpoint.Query()
	query.Set("long", fmt.Sprint(a.lon))
	query.Set("lat", fmt.Sprint(a.lat))
	query.Set("precision", "hourly")
	query.Set("date", normalisedTime(t).Format("2006-01-02"))
	query.Set("timezone", "UTC")
	query.Set("apiKey", a.apiKey)
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d: %s", resp.StatusCode, string(body))
	}

	var parsed apiFreaksHistoricalResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}
	if len(parsed.Historical.Hourly) == 0 {
		return fmt.Errorf("response contains no hourly weather data")
	}

	data := make(map[time.Time]*WeatherData, len(parsed.Historical.Hourly))
	for _, hour := range parsed.Historical.Hourly {
		timestamp, err := time.Parse("2006-01-02T15:04", hour.Timestamp)
		if err != nil {
			return fmt.Errorf("parse hourly timestamp %q: %w", hour.Timestamp, err)
		}
		data[timestamp] = &WeatherData{
			Code:          hour.Code,
			Temperature:   hour.Temperature,
			Precipitation: hour.Precipitation,
			WindSpeed:     hour.WindSpeed / 3.6,
		}
	}
	for timestamp, weather := range data {
		a.store.put(timestamp, weather)
	}
	return nil
}

type apiFreaksHistoricalResponse struct {
	Historical struct {
		Hourly []struct {
			Timestamp     string      `json:"timestamp"`
			Code          WeatherCode `json:"weather_code"`
			Temperature   float64     `json:"temperature_2m"`
			Precipitation float64     `json:"precipitation"`
			WindSpeed     float64     `json:"wind_speed_10m"`
		} `json:"hourly"`
	} `json:"historical"`
}

var _ HistoricalProvider = apiFreaks{}
