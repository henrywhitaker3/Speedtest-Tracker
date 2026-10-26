package notifications

import (
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
