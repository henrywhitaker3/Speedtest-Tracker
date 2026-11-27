package config

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/weather"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
)

func RegisterRoutes(
	h *http.HTTP,
	conf *config.Config,
	host *actor.Service,
	provider weather.Provider,
) {
	h.Register(&configHandler{conf: *conf})
	h.Register(&notificationsHandler{conf: conf.Notifications})
	h.Register(&listServersHandler{invoker: host})
	h.Register(&getCurrentWeatherHandler{p: provider})
	h.Register(&historicalWeatherStatusHandler{lister: host})
	h.Register(&historicalWeatherBackfillHandler{dispatcher: host})
}
