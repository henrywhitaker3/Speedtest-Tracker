package config

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/henrywhitaker3/windowframe/v2/config"
	"github.com/henrywhitaker3/windowframe/v2/crypto"
)

type NotificationPhase string

type NotificationKind string

const (
	Webhook     NotificationKind = "webhook"
	Healthcheck NotificationKind = "healthcheck"
	Telegram    NotificationKind = "telegram"
	Slack       NotificationKind = "slack"
	Discord     NotificationKind = "discord"

	PreRun  NotificationPhase = "pre-run"
	PostRun NotificationPhase = "post-run"
)

var (
	ValidNotificationKinds  = []NotificationKind{Webhook, Healthcheck, Telegram, Slack, Discord}
	ValidNotificationPhases = []NotificationPhase{PreRun, PostRun}
)

type Notification struct {
	Kind   NotificationKind    `json:"kind"`
	Phases []NotificationPhase `json:"phases"`

	URL string `json:"url,omitempty"`
	// Headers to add to webhook requests
	Headers map[string]string `json:"headers,omitempty"`

	// Telegram-specific options
	BotToken config.RedactedString `json:"botToken,omitempty"`
	ChatID   config.RedactedString `json:"chatID,omitempty"`
}

func (n Notification) String() string {
	by, _ := json.Marshal(n)
	sum := crypto.Sum(string(by))
	return hex.EncodeToString([]byte(sum))
}

func (n Notification) validate() error {
	if !slices.Contains(ValidNotificationKinds, n.Kind) {
		return fmt.Errorf("invalid kind %s", n.Kind)
	}
	for _, p := range n.Phases {
		if !slices.Contains(ValidNotificationPhases, p) {
			return fmt.Errorf("invalid phase %s", p)
		}
	}

	switch n.Kind {
	case Telegram:
		if n.BotToken == "" {
			return fmt.Errorf("telegram bot token must be set")
		}
		if n.ChatID == "" {
			return fmt.Errorf("telegram chat id must be set")
		}
	case Webhook:
		fallthrough
	case Discord:
		fallthrough
	case Slack:
		fallthrough
	case Healthcheck:
		if !strings.HasPrefix(n.URL, "https://") && !strings.HasPrefix(n.URL, "http://") {
			return fmt.Errorf("url does not start with http(s):// scheme")
		}
		if _, err := url.Parse(n.URL); err != nil {
			return fmt.Errorf("invalid url: %w", err)
		}
	default:
		return fmt.Errorf("unknown notification kind %s", n.Kind)
	}
	return nil
}

type Notifications []Notification

type notificationExtractor struct{}

func (notificationExtractor) Extract(conf *Config) error {
	for i := 0; ; i++ {
		prefix := fmt.Sprintf("NOTIFICATIONS_%d_", i)

		kind, hasKind := os.LookupEnv(prefix + "KIND")
		phasesRaw, hasPhases := os.LookupEnv(prefix + "PHASES")

		if !hasKind && !hasPhases {
			return nil
		}
		if !hasKind || !hasPhases {
			return fmt.Errorf("notification %d must have KIND, PHASES", i)
		}

		phases := []NotificationPhase{}
		for p := range strings.SplitSeq(phasesRaw, ",") {
			phases = append(phases, NotificationPhase(p))
		}

		var headers map[string]string
		if raw, ok := os.LookupEnv(prefix + "HEADERS"); ok {
			if err := json.Unmarshal([]byte(raw), &headers); err != nil {
				return fmt.Errorf("unmarshal headers: %w", err)
			}
		}

		notification := Notification{
			Kind:     NotificationKind(kind),
			Phases:   phases,
			URL:      os.Getenv(prefix + "URL"),
			Headers:  headers,
			BotToken: config.RedactedString(os.Getenv("NOTIFICATIONS_TELEGRAM_BOT_TOKEN")),
			ChatID:   config.RedactedString(os.Getenv("NOTIFICATIONS_TELEGRAM_CHAT_ID")),
		}
		if err := notification.validate(); err != nil {
			return fmt.Errorf("notification %d: %w", i, err)
		}
		conf.Notifications = append(conf.Notifications, notification)
	}
}
