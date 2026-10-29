package notifications

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
)

type NotifierOpts struct {
	Host          *local.Host
	UIURL         string
	Notifications config.Notifications
	Telegram      telegramMessageSender
}

type Notifier struct {
	client        actor.Client[notifierState]
	uiURL         string
	notifications config.Notifications
	telegram      telegramMessageSender
}

func RegisterActors(opts NotifierOpts) error {
	if err := opts.Host.RegisterActor(Actor, newNotifier(opts)); err != nil {
		return fmt.Errorf("register notifier actor: %w", err)
	}
	return nil
}

type notifierState struct {
	Result *results.Result
	Record map[string]bool
}

const (
	Actor = "notifications"
)

func newNotifier(opts NotifierOpts) actor.Factory {
	return func(actorID string, service *actor.Service) actor.Actor {
		return &Notifier{
			client:        actor.NewActorClient[notifierState](Actor, actorID, service),
			uiURL:         opts.UIURL,
			notifications: opts.Notifications,
			telegram:      opts.Telegram,
		}
	}
}

const (
	Notify       = "notify"
	ShouldNotify = "should-notify"
	stateTTL     = time.Minute * 5
)

func (n *Notifier) Invoke(ctx context.Context, method string, data actor.Envelope) (any, error) {
	switch method {
	case Notify:
		return n.invokeNotify(ctx, data)
	case ShouldNotify:
		return n.shouldNotify(), nil
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func (n *Notifier) shouldNotify() bool {
	return len(n.notifications) > 0
}

type NotifyInput struct {
	Result results.Result
	Phase  Phase
}

func (n *Notifier) invokeNotify(
	ctx context.Context,
	data actor.Envelope,
) (*NotificationError, error) {
	state, err := n.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get notifier state: %w", err)
	}

	var input NotifyInput
	if err := data.Decode(&input); err != nil {
		return nil, fmt.Errorf("decode result input: %w", err)
	}

	if state.Result != nil {
		if state.Result.ID != input.Result.ID {
			return nil, fmt.Errorf("incorrect result passed to actor")
		}
	}
	state.Result = &input.Result
	if state.Record == nil {
		state.Record = map[string]bool{}
	}
	if err := n.client.SetState(ctx, state, &actor.SetStateOpts{
		TTL: stateTTL,
	}); err != nil {
		return nil, fmt.Errorf("store state: %w", err)
	}

	errors := &NotificationError{}
	for _, no := range n.notifications {
		if passed, exists := state.Record[fmt.Sprintf("%d:%s", input.Phase, no.String())]; passed &&
			exists {
			continue
		}
		if err := notify(ctx, n.uiURL, n.telegram, no, input.Phase, input.Result); err != nil {
			errors.push(no, err)
		} else {
			state.Record[fmt.Sprintf("%d:%s", input.Phase, no.String())] = true
		}
	}

	if err := n.client.SetState(ctx, state, &actor.SetStateOpts{
		TTL: stateTTL,
	}); err != nil {
		return nil, fmt.Errorf("store state: %w", err)
	}

	if errors.Nil() {
		return nil, nil
	}

	return errors, nil
}

func notify(
	ctx context.Context,
	uiURL string,
	telegram telegramMessageSender,
	n config.Notification,
	phase Phase,
	result results.Result,
) error {
	nctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var err error
	switch n.Kind {
	case config.Healthcheck:
		err = notifyHealthcheck(nctx, n.URL, phase, result)
	case config.Webhook:
		err = sendWebhook(nctx, n.URL, result)
	case config.Telegram:
		err = SendTelegram(nctx, telegram, n, uiURL, phase, result)
	case config.Slack:
		err = sendSlack(nctx, n.URL, uiURL, phase, result)
	case config.Discord:
		err = sendDiscord(nctx, n.URL, uiURL, phase, result)
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("send %s: %w", n.Kind, err)
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
