package weather

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItMarshalsWeatherCodes(t *testing.T) {
	tcs := []struct {
		code     WeatherCode
		expected string
	}{
		{
			code:     0,
			expected: `{"emoji":"☀️","name":"Clear sky"}`,
		},
	}

	for _, c := range tcs {
		t.Run(fmt.Sprint(c.code), func(t *testing.T) {
			by, err := json.Marshal(c.code)
			require.Nil(t, err)
			require.Equal(t, c.expected, string(by))
		})
	}
}
