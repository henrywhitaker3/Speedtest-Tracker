// Package api
package api

import (
	"log/slog"

	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
)

func New(url, version string, service *actor.Service) *http.HTTP {
	h := http.New(http.HTTPOpts{
		Port: 8080,
		Openapi: http.OpenapiOpts{
			Enabled:        true,
			ServiceName:    "Speedtest-Tracker",
			ServiceVersion: version,
			PublicURL:      url,
		},
		Logger: slog.With("component", "api"),
	})

	registerRoutes(h, service)

	return h
}

func registerRoutes(h *http.HTTP, service *actor.Service) {
	h.Register(&runHandler{service: service})
}
