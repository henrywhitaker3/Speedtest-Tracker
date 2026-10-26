package notifications

import (
	"slices"
	"strings"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
)

type NotificationError struct {
	errors     []error
	redactions []string
}

func (ne *NotificationError) push(n config.Notification, err error) {
	ne.errors = append(ne.errors, err)
	if n.URL != "" {
		ne.redactions = append(ne.redactions, n.URL)
		if n.BotToken != "" {
			if !slices.Contains(ne.redactions, string(n.BotToken)) {
				ne.redactions = append(ne.redactions, string(n.BotToken))
			}
		}
		if n.ChatID != "" {
			if !slices.Contains(ne.redactions, string(n.ChatID)) {
				ne.redactions = append(ne.redactions, string(n.ChatID))
			}
		}
	}
}

func (ne NotificationError) Error() string {
	var out []string
	for _, e := range ne.errors {
		str := e.Error()
		for _, r := range ne.redactions {
			str = strings.ReplaceAll(str, r, "[REDACTED]")
		}
		out = append(out, str)
	}
	return strings.Join(out, ", ")
}
