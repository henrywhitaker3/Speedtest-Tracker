package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotify/go-api-client/v2/models"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestSendGotify(t *testing.T) {
	ping := 7.5
	download := int64(10_000_000)
	upload := int64(5_000_000)
	result := results.Result{
		ID:       uuid.MustOrdered(),
		Ping:     &ping,
		Download: &download,
		Upload:   &upload,
	}
	failedResult := result
	failure := "speedtest failed"
	failedResult.Error = &failure

	tests := []struct {
		name         string
		phase        Phase
		result       results.Result
		wantTitle    string
		wantMessage  string
		wantPriority int
	}{
		{
			name:         "sends a starting message before a speedtest",
			phase:        PreRun,
			result:       result,
			wantTitle:    StartingTitle,
			wantMessage:  "[See result](https://speedtest.example/results/" + result.ID.String() + ")",
			wantPriority: 5,
		},
		{
			name:         "sends measurements after a speedtest",
			phase:        PostRun,
			result:       result,
			wantTitle:    SuccessTitle,
			wantMessage:  "Ping: **7.50 ms**\nDownload: **80.00 Mbps**\nUpload: **40.00 Mbps**\n[See result](https://speedtest.example/results/" + result.ID.String() + ")",
			wantPriority: 5,
		},
		{
			name:         "sends the failure reason after a failed speedtest",
			phase:        PostRun,
			result:       failedResult,
			wantTitle:    ErrorTitle,
			wantMessage:  "Error: speedtest failed\n[See result](https://speedtest.example/results/" + result.ID.String() + ")",
			wantPriority: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var payload models.MessageExternal
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					require.Equal(t, "/message", r.URL.Path)
					require.Equal(t, "application-token", r.Header.Get("X-Gotify-Key"))
					require.Equal(t, "application/json", r.Header.Get("Content-Type"))
					require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
					w.Header().Set("Content-Type", "application/json")
					_, err := w.Write([]byte(`{}`))
					require.NoError(t, err)
				}),
			)
			t.Cleanup(server.Close)

			err := sendGotify(
				t.Context(),
				server.URL,
				"application-token",
				"https://speedtest.example",
				tt.phase,
				tt.result,
			)
			require.NoError(t, err)
			require.Equal(t, tt.wantTitle, payload.Title)
			require.Equal(t, tt.wantMessage, payload.Message)
			require.Equal(t, tt.wantPriority, payload.Priority)
			require.Equal(t, map[string]any{
				"client::display": map[string]any{"contentType": "text/markdown"},
				"client::notification": map[string]any{
					"click": map[string]any{
						"url": "https://speedtest.example/results/" + result.ID.String(),
					},
				},
			}, payload.Extras)
		})
	}
}

func TestBuildGotifyMessageRejectsUnknownPhase(t *testing.T) {
	message, err := buildGotifyMessage("https://speedtest.example", Phase(99), results.Result{})
	require.Nil(t, message)
	require.EqualError(t, err, "unknown phase 99")
}
