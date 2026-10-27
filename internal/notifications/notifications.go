// Package notifications
package notifications

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

type NotifierOpts struct {
	UIURL         string
	Notifications config.Notifications
	Telegram      telegramMessageSender
}

type Notifier struct {
	uiURL         string
	notifications config.Notifications
	telegram      telegramMessageSender
}

func New(opts NotifierOpts) *Notifier {
	return &Notifier{
		notifications: opts.Notifications,
		uiURL:         opts.UIURL,
		telegram:      opts.Telegram,
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
		case config.Telegram:
			if no.telegram == nil {
				errors.push(n, fmt.Errorf("telegram sender is not configured"))
			} else if err := SendTelegram(nctx, no.telegram, n, no.uiURL, phase, result); err != nil {
				errors.push(n, err)
			}
		case config.Slack:
			if err := sendSlack(ctx, n.URL, no.uiURL, phase, result); err != nil {
				errors.push(n, err)
			}
		case config.Discord:
			if err := sendDiscord(ctx, n.URL, no.uiURL, phase, result); err != nil {
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
