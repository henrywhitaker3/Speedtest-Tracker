package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	openWeatherMapBaseURL = "https://api.openweathermap.org"
)

type openWeatherMap struct {
	apiKey   string
	lon, lat float64
	baseURL  string
}

func newOpenWeatherMap(apiKey string, lon, lat float64) openWeatherMap {
	return openWeatherMap{
		apiKey:  apiKey,
		lon:     lon,
		lat:     lat,
		baseURL: openWeatherMapBaseURL,
	}
}

func (o openWeatherMap) Get(ctx context.Context) (*WeatherData, error) {
	endpoint, err := url.Parse(o.baseURL + "/data/2.5/weather")
	if err != nil {
		return nil, fmt.Errorf("build url: %w", err)
	}
	query := endpoint.Query()
	query.Set("lon", fmt.Sprint(o.lon))
	query.Set("lat", fmt.Sprint(o.lat))
	query.Set("appid", o.apiKey)
	query.Set("units", "metric")
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code %d: %s", resp.StatusCode, string(body))
	}

	var parsed openWeatherMapData
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("parse response body: %w", err)
	}

	if len(parsed.Weather) == 0 {
		return nil, fmt.Errorf("response has no weather conditions")
	}

	return &WeatherData{
		Code:          openWeatherMapCode(parsed.Weather[0].ID),
		Temperature:   parsed.Main.Temperature,
		Precipitation: parsed.Rain.OneHour + parsed.Snow.OneHour,
		WindSpeed:     parsed.Wind.Speed,
	}, nil
}

type openWeatherMapData struct {
	Weather []struct {
		ID int `json:"id"`
	} `json:"weather"`
	Main struct {
		Temperature float64 `json:"temp"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Rain struct {
		OneHour float64 `json:"1h"`
	} `json:"rain"`
	Snow struct {
		OneHour float64 `json:"1h"`
	} `json:"snow"`
}

func openWeatherMapCode(id int) WeatherCode {
	switch id {
	case 200, 201, 210, 211, 221, 230, 231:
		return 95
	case 202, 212, 232:
		return 97
	case 300, 310:
		return 51
	case 301, 311, 313, 321:
		return 53
	case 302, 312, 314:
		return 55
	case 500:
		return 61
	case 501:
		return 63
	case 502, 503, 504:
		return 65
	case 511:
		return 66
	case 520:
		return 80
	case 521:
		return 81
	case 522, 531:
		return 82
	case 600, 611, 612, 615:
		return 71
	case 601, 613, 616:
		return 73
	case 602:
		return 75
	case 620, 621:
		return 85
	case 622:
		return 86
	case 701, 741:
		return 45
	case 800:
		return 0
	case 801:
		return 1
	case 802, 803:
		return 2
	case 804:
		return 3
	default:
		return -1
	}
}

var _ Provider = openWeatherMap{}
