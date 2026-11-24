package config

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/francis"
	"github.com/henrywhitaker3/windowframe/v2/http"
)

func RegisterRoutes(h *http.HTTP, conf *config.Config, client francis.Invoker) {
	h.Register(&configHandler{conf: *conf})
	h.Register(&notificationsHandler{conf: conf.Notifications})
	h.Register(&listServersHandler{invoker: client})
}
