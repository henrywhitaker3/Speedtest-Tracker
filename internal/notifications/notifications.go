// Package notifications
package notifications

import (
	"context"
	"slices"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
)

type telegramSender func(
	ctx context.Context,
	botToken string,
	chatID string,
	uiURL string,
	phase Phase,
	result results.Result,
) error

type NotifierOpts struct {
	UIURL         string
	Notifications config.Notifications
	Telegram      telegramSender
}

type Notifier struct {
	uiURL         string
	notifications config.Notifications
	telegram      telegramSender
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
