package weather

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/swaggest/jsonschema-go"
)

func TestWeatherCodeSchema(t *testing.T) {
	reflector := jsonschema.Reflector{}
	schema, err := reflector.Reflect(WeatherData{})
	require.NoError(t, err)
	encoded, err := json.Marshal(schema)
	require.NoError(t, err)
	var document map[string]any
	require.NoError(t, json.Unmarshal(encoded, &document))
	properties := document["properties"].(map[string]any)
	require.Equal(t, "#/definitions/WeatherWeatherCode", properties["code"].(map[string]any)["$ref"])
	code := document["definitions"].(map[string]any)["WeatherWeatherCode"].(map[string]any)
	require.Equal(t, "object", code["type"])
	require.ElementsMatch(t, []any{"name", "emoji"}, code["required"])
	codeProperties := code["properties"].(map[string]any)
	require.Equal(t, "string", codeProperties["name"].(map[string]any)["type"])
	require.Equal(t, "string", codeProperties["emoji"].(map[string]any)["type"])
}

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
