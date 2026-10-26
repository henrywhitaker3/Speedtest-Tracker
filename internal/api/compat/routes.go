// Package compat
package compat

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"gorm.io/gorm"
)

func Register(h *http.HTTP, db *gorm.DB) {
	h.Register(&latestHandler{repo: results.NewRepo(db)})
}
