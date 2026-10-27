package notifications

import (
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/slack-go/slack"
	"github.com/stretchr/testify/require"
)

func TestBuildSlackMessages(t *testing.T) {
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
	resultURL := "https://speedtest.example/results/" + result.ID.String()

	tests := []struct {
		name           string
		message        *slack.WebhookMessage
		wantText       string
		wantHeader     string
		wantFields     []string
		wantBody       string
		wantBlockCount int
	}{
		{
			name:           "builds a starting message",
			message:        buildSlackPreRun(result, "https://speedtest.example"),
			wantText:       "Starting Speedtest",
			wantHeader:     "Starting Speedtest",
			wantBlockCount: 2,
		},
		{
			name:       "builds a completed message with measurements",
			message:    buildSlackPostRun(result, "https://speedtest.example"),
			wantText:   "Finished Speedtest",
			wantHeader: "Finished Speedtest",
			wantFields: []string{
				"*Ping*\n7.50 ms",
				"*Download*\n80.00 Mbps",
				"*Upload*\n40.00 Mbps",
			},
			wantBlockCount: 3,
		},
		{
			name:           "builds a failed message with the error",
			message:        buildSlackPostRun(failedResult, "https://speedtest.example"),
			wantText:       "Failed Speedtest",
			wantHeader:     "Failed Speedtest",
			wantBody:       "*Error*\nspeedtest failed",
			wantBlockCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantText, tt.message.Text)
			require.NotNil(t, tt.message.Blocks)
			require.Len(t, tt.message.Blocks.BlockSet, tt.wantBlockCount)

			header, ok := tt.message.Blocks.BlockSet[0].(*slack.HeaderBlock)
			require.True(t, ok)
			require.Equal(t, tt.wantHeader, header.Text.Text)

			resultLink, ok := tt.message.Blocks.BlockSet[tt.wantBlockCount-1].(*slack.SectionBlock)
			require.True(t, ok)
			require.Equal(t, "<"+resultURL+"|View Result>", resultLink.Text.Text)

			if tt.wantFields != nil {
				measurements, ok := tt.message.Blocks.BlockSet[1].(*slack.SectionBlock)
				require.True(t, ok)
				require.Len(t, measurements.Fields, len(tt.wantFields))
				for index, wantField := range tt.wantFields {
					require.Equal(t, wantField, measurements.Fields[index].Text)
				}
			}

			if tt.wantBody != "" {
				body, ok := tt.message.Blocks.BlockSet[1].(*slack.SectionBlock)
				require.True(t, ok)
				require.Equal(t, tt.wantBody, body.Text.Text)
			}
		})
	}
}
