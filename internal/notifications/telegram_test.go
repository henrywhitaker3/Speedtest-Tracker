package notifications

import (
	"fmt"
	"testing"

	"github.com/go-telegram/bot/models"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/stretchr/testify/require"
)

func TestTelegramMessageBuildersUseMarkdownResultLinks(t *testing.T) {
	result := results.Result{ID: uuid.MustOrdered()}
	resultURL := "https://speedtest.example/results/" + result.ID.String()

	require.Equal(t, resultURL, buildResultsURL("https://speedtest.example", result))
	require.Equal(
		t,
		"[See result]("+resultURL+")",
		buildMarkdownURL("https://speedtest.example", result),
	)
}

func TestNotifierSendsTelegramMessages(t *testing.T) {
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
	sendErr := fmt.Errorf("telegram unavailable")

	tests := []struct {
		name    string
		phase   Phase
		result  results.Result
		sendErr error
		want    string
		wantErr bool
	}{
		{
			name:   "sends a starting message before a speedtest",
			phase:  PreRun,
			result: result,
			want:   fmt.Sprintf("*%s*\n[See result](%s)", StartingTitle, resultURL),
		},
		{
			name:   "sends completed measurements after a speedtest",
			phase:  PostRun,
			result: result,
			want: fmt.Sprintf(
				"*%s*\nPing: *7.50 ms*\nDownload: *80.00 Mbps*\nUpload: *40.00 Mbps*\n[See result](%s)",
				SuccessTitle,
				resultURL,
			),
		},
		{
			name:   "sends an error after a failed speedtest",
			phase:  PostRun,
			result: failedResult,
			want: fmt.Sprintf(
				"*%s*\nError: %s\n[See result](%s)",
				ErrorTitle,
				failure,
				resultURL,
			),
		},
		{
			name:    "returns a delivery error",
			phase:   PreRun,
			result:  result,
			sendErr: sendErr,
			want:    fmt.Sprintf("*%s*\n[See result](%s)", StartingTitle, resultURL),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender := &telegramSenderStub{err: tt.sendErr}
			notifier := New(NotifierOpts{
				UIURL: "https://speedtest.example",
				Notifications: config.Notifications{{
					Kind:   config.Telegram,
					Phases: []config.NotificationPhase{notificationPhase(tt.phase)},
					ChatID: "123456",
				}},
				Telegram: sender,
			})

			err := notifier.Notify(t.Context(), tt.phase, tt.result)
			if tt.wantErr {
				require.ErrorContains(t, err, sendErr.Error())
			} else {
				require.NoError(t, err)
			}

			require.Len(t, sender.params, 1)
			params := sender.params[0]
			require.Equal(t, "123456", params.ChatID)
			require.Equal(t, models.ParseModeMarkdownV1, params.ParseMode)
			require.Equal(t, tt.want, params.Text)
		})
	}
}
