// Package speedtest
package speedtest

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/actor"
	"gorm.io/gorm"
)

func Register(h *http.HTTP, db *gorm.DB, service *actor.Service) {
	repo := results.NewRepo(db)
	h.Register(&runHandler{service: service})
	h.Register(&getHandler{repo: repo})
	h.Register(&listHandler{repo: repo})
}
