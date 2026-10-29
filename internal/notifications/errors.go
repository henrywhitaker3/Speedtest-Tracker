package notifications

import (
	"slices"
	"strings"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
)

type NotificationError struct {
	Errors     []error
	Redactions []string
}

func (n NotificationError) Nil() bool {
	return len(n.Errors) == 0
}

func (ne *NotificationError) push(n config.Notification, err error) {
	ne.Errors = append(ne.Errors, err)
	if n.URL != "" {
		ne.Redactions = append(ne.Redactions, n.URL)
		withoutScheme := strings.ReplaceAll(n.URL, "http://", "")
		withoutScheme = strings.ReplaceAll(withoutScheme, "https://", "")
		ne.Redactions = append(ne.Redactions, withoutScheme)
		if n.BotToken != "" {
			if !slices.Contains(ne.Redactions, string(n.BotToken)) {
				ne.Redactions = append(ne.Redactions, string(n.BotToken))
			}
		}
		if n.ChatID != "" {
			if !slices.Contains(ne.Redactions, string(n.ChatID)) {
				ne.Redactions = append(ne.Redactions, string(n.ChatID))
			}
		}
	}
}

func (ne NotificationError) Error() string {
	var out []string
	for _, e := range ne.Errors {
		str := e.Error()
		for _, r := range ne.Redactions {
			str = strings.ReplaceAll(str, r, "[REDACTED]")
		}
		out = append(out, str)
	}
	return strings.Join(out, ", ")
}
