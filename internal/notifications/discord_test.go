package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	discord "github.com/bwmarrin/discordgo"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestSendDiscord(t *testing.T) {
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
		name       string
		phase      Phase
		result     results.Result
		wantTitle  string
		wantColor  int
		wantFields []*discord.MessageEmbedField
	}{
		{
			name:      "sends a starting embed before a speedtest",
			phase:     PreRun,
			result:    result,
			wantTitle: "Starting Speedtest",
			wantColor: discordInfoColor,
		},
		{
			name:      "sends measurements after a speedtest",
			phase:     PostRun,
			result:    result,
			wantTitle: "Finished Speedtest",
			wantColor: discordSuccessColor,
			wantFields: []*discord.MessageEmbedField{
				{Name: "Ping", Value: "7.50 ms", Inline: true},
				{Name: "Download", Value: "80.00 Mbps", Inline: true},
				{Name: "Upload", Value: "40.00 Mbps", Inline: true},
			},
		},
		{
			name:      "sends the failure reason after a failed speedtest",
			phase:     PostRun,
			result:    failedResult,
			wantTitle: "Speedtest Failed",
			wantColor: discordFailureColor,
			wantFields: []*discord.MessageEmbedField{
				{Name: "Error", Value: failure},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &discord.WebhookParams{}
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					require.Equal(t, http.MethodPost, r.Method)
					require.Equal(t, "application/json", r.Header.Get("Content-Type"))
					require.NoError(t, json.NewDecoder(r.Body).Decode(params))
				}),
			)
			t.Cleanup(server.Close)

			err := sendDiscord(
				t.Context(),
				server.URL,
				"https://speedtest.example",
				tt.phase,
				tt.result,
			)
			require.NoError(t, err)
			require.Len(t, params.Embeds, 1)
			require.Equal(t, tt.wantTitle, params.Embeds[0].Title)
			require.Equal(
				t,
				"https://speedtest.example/results/"+tt.result.ID.String(),
				params.Embeds[0].URL,
			)
			require.Equal(t, tt.wantColor, params.Embeds[0].Color)
			require.Equal(t, tt.wantFields, params.Embeds[0].Fields)
		})
	}
}
