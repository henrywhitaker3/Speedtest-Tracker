// Package notifications
package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

type NotificationError struct {
	errors     []error
	redactions []string
}

func (ne *NotificationError) push(n config.Notification, err error) {
	ne.errors = append(ne.errors, err)
	ne.redactions = append(ne.redactions, n.URL)
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

type Notifier struct {
	notifications config.Notifications
}

func New(n config.Notifications) *Notifier {
	return &Notifier{
		notifications: n,
	}
}

type Phase int

const (
	PreRun Phase = iota
	PostRun
)

func (no Notifier) ShouldNotify() bool {
	return len(no.notifications) > 0
}

func (no Notifier) Notify(ctx context.Context, phase Phase, result results.Result) error {
	errors := &NotificationError{}

	for _, n := range no.notifications {
		if !notifiesAtPhase(n, phase) {
			continue
		}

		nctx, cancel := context.WithTimeout(ctx, time.Second*5)
		switch n.Kind {
		case config.Healthcheck:
			if err := notifyHealthcheck(nctx, n.URL, phase, result); err != nil {
				errors.push(n, err)
			}
		case config.Webhook:
			if err := sendWebhook(nctx, n.URL, result); err != nil {
				errors.push(n, err)
			}
		}
		cancel()
	}

	if len(errors.errors) == 0 {
		return nil
	}

	return errors
}

func notifiesAtPhase(notification config.Notification, phase Phase) bool {
	switch phase {
	case PreRun:
		return slices.Contains(notification.Phases, config.PreRun)
	case PostRun:
		return slices.Contains(notification.Phases, config.PostRun)
	default:
		return false
	}
}

func notifyHealthcheck(ctx context.Context, url string, phase Phase, result results.Result) error {
	switch phase {
	case PreRun:
		url += "/start"
	case PostRun:
		if result.Error != nil {
			url += "/fail"
		}
	}
	return sendWebhook(ctx, url, result)
}

func sendWebhook(ctx context.Context, url string, result results.Result) error {
	body, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	bodyReader := bytes.NewReader(body)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bodyReader,
	)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do healthcheck request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}

	return nil
}
