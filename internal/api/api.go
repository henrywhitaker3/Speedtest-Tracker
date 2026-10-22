// Package api
package api

import (
	"log/slog"

	"github.com/henrywhitaker3/Speedtest-Tracker/internal/api/speedtest"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
	"gorm.io/gorm"
)

func New(url, version string, db *gorm.DB, service *actor.Service) *http.HTTP {
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

	registerRoutes(h, db, service)

	return h
}

func registerRoutes(h *http.HTTP, db *gorm.DB, service *actor.Service) {
	speedtest.Register(h, db, service)
}
