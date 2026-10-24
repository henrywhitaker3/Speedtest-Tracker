package config

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"
)

type NotificationPhase string

type NotificationKind string

const (
	Webhook     NotificationKind = "webhook"
	Healthcheck NotificationKind = "healthcheck"

	PreRun  NotificationPhase = "pre-run"
	PostRun NotificationPhase = "post-run"
)

var (
	ValidNotificationKinds  = []NotificationKind{Webhook, Healthcheck}
	ValidNotificationPhases = []NotificationPhase{PreRun, PostRun}
)

type Notification struct {
	Kind   NotificationKind
	Phases []NotificationPhase
	URL    string
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
	if !strings.HasPrefix(n.URL, "https://") && !strings.HasPrefix("http://", n.URL) {
		return fmt.Errorf("url does not start with http(s):// scheme")
	}
	if _, err := url.Parse(n.URL); err != nil {
		return fmt.Errorf("invalid url: %w", err)
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
		url, hasURL := os.LookupEnv(prefix + "URL")

		if !hasKind && !hasPhases && !hasURL {
			return nil
		}
		if !hasKind || !hasPhases || !hasURL {
			return fmt.Errorf("notification %d must have KIND, PHASES and URL", i)
		}

		phases := []NotificationPhase{}
		for _, p := range strings.Split(phasesRaw, ",") {
			phases = append(phases, NotificationPhase(p))
		}

		notification := Notification{
			Kind:   NotificationKind(kind),
			Phases: phases,
			URL:    url,
		}
		if err := notification.validate(); err != nil {
			return fmt.Errorf("notification %d: %w", i, err)
		}
		conf.Notifications = append(conf.Notifications, notification)
	}
}
