// Package auth
package auth

import (
	"fmt"
	"strings"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/host/local"
)

func RegisterRoutes(h *http.HTTP, host *local.Host, conf config.Auth, appURL string) error {
	if !conf.Enabled {
		return nil
	}

	if appURL == "" {
		return fmt.Errorf("build OIDC callback URL: app URL is required")
	}
	callbackURL := strings.TrimRight(appURL, "/") + "/v1/auth/callback"

	h.Register(&loginHandler{
		conf:        conf,
		callbackURL: callbackURL,
	})
	h.Register(&callbackHandler{service: host.Service()})
	return nil
}
