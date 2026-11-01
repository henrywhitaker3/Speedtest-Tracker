package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItParsesNotifications(t *testing.T) {
	tcs := []struct {
		name  string
		env   map[string]string
		want  Notifications
		error bool
	}{
		{
			name: "empty slice with no env vars set",
		},
		{
			name: "errors with kind not set",
			env: map[string]string{
				"NOTIFICATIONS_0_PHASES": "pre-run",
				"NOTIFICATIONS_0_URL":    "https://example.com/healthcheck",
			},
			error: true,
		},
		{
			name: "errors with url not set",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "healthcheck",
				"NOTIFICATIONS_0_PHASES": "pre-run",
			},
			error: true,
		},
		{
			name: "errors with phases not set",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND": "healthcheck",
				"NOTIFICATIONS_0_URL":  "https://example.com/healthcheck",
			},
			error: true,
		},
		{
			name: "single slice when one is set",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "healthcheck",
				"NOTIFICATIONS_0_PHASES": "pre-run,post-run",
				"NOTIFICATIONS_0_URL":    "https://example.com/healthcheck",
			},
			want: Notifications{{
				Kind:   Healthcheck,
				Phases: []NotificationPhase{PreRun, PostRun},
				URL:    "https://example.com/healthcheck",
			}},
		},
		{
			name: "parses a gotify notification",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":         "gotify",
				"NOTIFICATIONS_0_PHASES":       "pre-run,post-run",
				"NOTIFICATIONS_0_URL":          "https://gotify.example",
				"NOTIFICATIONS_0_GOTIFY_TOKEN": "application-token",
			},
			want: Notifications{{
				Kind:        Gotify,
				Phases:      []NotificationPhase{PreRun, PostRun},
				URL:         "https://gotify.example",
				GotifyToken: "application-token",
			}},
		},
		{
			name: "errors when a gotify notification has no token",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "gotify",
				"NOTIFICATIONS_0_PHASES": "pre-run",
				"NOTIFICATIONS_0_URL":    "https://gotify.example",
			},
			error: true,
		},
		{
			name: "multiple notifications parsed",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "webhook",
				"NOTIFICATIONS_0_PHASES": "post-run",
				"NOTIFICATIONS_0_URL":    "https://example.com/webhook",
				"NOTIFICATIONS_1_KIND":   "healthcheck",
				"NOTIFICATIONS_1_PHASES": "pre-run,post-run",
				"NOTIFICATIONS_1_URL":    "https://example.com/healthcheck",
			},
			want: Notifications{
				{
					Kind:   Webhook,
					Phases: []NotificationPhase{PostRun},
					URL:    "https://example.com/webhook",
				},
				{
					Kind:   Healthcheck,
					Phases: []NotificationPhase{PreRun, PostRun},
					URL:    "https://example.com/healthcheck",
				},
			},
		},
		{
			name: "healthcheck env var produces a notification item",
			env: map[string]string{
				"NOTIFICATIONS_HEALTHCHECKS": "https://example.com/healthcheck",
			},
			want: Notifications{{
				Kind:   Healthcheck,
				Phases: []NotificationPhase{PreRun, PostRun},
				URL:    "https://example.com/healthcheck",
			}},
		},
		{
			name: "errors with invalid kind",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "email",
				"NOTIFICATIONS_0_PHASES": "pre-run",
				"NOTIFICATIONS_0_URL":    "https://example.com/notification",
			},
			error: true,
		},
		{
			name: "errors with invalid phase",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "webhook",
				"NOTIFICATIONS_0_PHASES": "pre-flight",
				"NOTIFICATIONS_0_URL":    "https://example.com/notification",
			},
			error: true,
		},
		{
			name: "errors with invalid url",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "healthcheck",
				"NOTIFICATIONS_0_PHASES": "post-run",
				"NOTIFICATIONS_0_URL":    "https://example.com/%zz",
			},
			error: true,
		},
		{
			name: "allows http urls",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "healthcheck",
				"NOTIFICATIONS_0_PHASES": "post-run",
				"NOTIFICATIONS_0_URL":    "http://example.com/",
			},
			want: Notifications{
				{
					Kind:   Healthcheck,
					Phases: []NotificationPhase{PostRun},
					URL:    "http://example.com/",
				},
			},
		},
		{
			name: "errors with non http(s) url",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":   "healthcheck",
				"NOTIFICATIONS_0_PHASES": "post-run",
				"NOTIFICATIONS_0_URL":    "ftp://example.com/",
			},
			error: true,
		},
		{
			name: "parses correctly formatted webhook headers",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":    "healthcheck",
				"NOTIFICATIONS_0_PHASES":  "post-run",
				"NOTIFICATIONS_0_URL":     "http://example.com/",
				"NOTIFICATIONS_0_HEADERS": `{"X-Bongo": "bingo"}`,
			},
			want: Notifications{
				{
					Kind:   Healthcheck,
					Phases: []NotificationPhase{PostRun},
					URL:    "http://example.com/",
					Headers: map[string]string{
						"X-Bongo": "bingo",
					},
				},
			},
		},
		{
			name: "errors with invalid format for headers",
			env: map[string]string{
				"NOTIFICATIONS_0_KIND":    "healthcheck",
				"NOTIFICATIONS_0_PHASES":  "post-run",
				"NOTIFICATIONS_0_URL":     "http://example.com/",
				"NOTIFICATIONS_0_HEADERS": `"bongo": "bingo"`,
			},
			error: true,
		},
	}

	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("ACCEPT_EULA", "true")
			for key, value := range c.env {
				t.Setenv(key, value)
			}

			conf, err := Load()
			if c.error {
				require.NotNil(t, err)
				return
			}
			require.Nil(t, err)
			require.Equal(t, c.want, conf.Notifications)
		})
	}
}

func TestItParsesServerIDs(t *testing.T) {
	tcs := []struct {
		name     string
		val      string
		expected []int
	}{
		{
			name:     "single server id",
			val:      "1",
			expected: []int{1},
		},
		{
			name:     "multiple comma-separated",
			val:      "1,2,3",
			expected: []int{1, 2, 3},
		},
	}
	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("ACCEPT_EULA", "true")
			t.Setenv("TEST_SERVER", c.val)
			conf, err := Load()
			require.Nil(t, err)
			require.Equal(t, c.expected, conf.Server)
		})
	}
}
