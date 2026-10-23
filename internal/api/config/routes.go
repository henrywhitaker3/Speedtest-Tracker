package config

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http"
)

func RegisterRoutes(h *http.HTTP, conf *config.Config) {
	h.Register(&configHandler{conf: *conf})
}
