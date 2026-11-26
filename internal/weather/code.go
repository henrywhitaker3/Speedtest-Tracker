package weather

import "encoding/json"

type WeatherCode int

func (w WeatherCode) MarshalJSON() ([]byte, error) {
	data := map[string]string{
		"name":  w.String(),
		"emoji": w.Emoji(),
	}
	return json.Marshal(data)
}

func (w WeatherCode) String() string {
	switch w {
	case 0:
		return "Clear sky"
	case 1:
		return "Mainly clear"
	case 2:
		return "Partly cloudy"
	case 3:
		return "Overcast"
	case 45:
		return "Fog"
	case 48:
		return "Depositing rime fog"
	case 51:
		return "Light drizzle"
	case 53:
		return "Moderate drizzle"
	case 55:
		return "Dense drizzle"
	case 56:
		return "Light freezing drizzle"
	case 57:
		return "Dense freezing drizzle"
	case 61:
		return "Slight rain"
	case 63:
		return "Moderate rain"
	case 65:
		return "Heavy rain"
	case 66:
		return "Light freezing rain"
	case 67:
		return "Heavy freezing rain"
	case 71:
		return "Slight snowfall"
	case 73:
		return "Moderate snowfall"
	case 75:
		return "Heavy snowfall"
	case 77:
		return "Snow grains"
	case 80:
		return "Slight rain showers"
	case 81:
		return "Moderate rain showers"
	case 82:
		return "Violent rain showers"
	case 85:
		return "Slight snow showers"
	case 86:
		return "Heavy snow showers"
	case 95:
		return "Thunderstorm"
	case 96:
		return "Thunderstorm with slight hail"
	case 97:
		return "Heavy thunderstorm"
	case 99:
		return "Thunderstorm with heavy hail"
	default:
		return "Unknown"
	}
}

func (w WeatherCode) Emoji() string {
	switch w {
	case 0:
		return "☀️"
	case 1:
		return "🌤️"
	case 2:
		return "⛅"
	case 3:
		return "☁️"
	case 45, 48:
		return "🌫️"
	case 51, 53, 55, 56, 57:
		return "🌦️"
	case 61, 63, 65, 66, 67, 80, 81, 82:
		return "🌧️"
	case 71, 73, 75, 77, 85, 86:
		return "🌨️"
	case 95, 96, 97, 99:
		return "⛈️"
	default:
		return "❓"
	}
}
