// Package auth
package auth

import (
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/host/local"
	"gorm.io/gorm"
)

func RegisterRoutes(
	h *http.HTTP,
	host *local.Host,
	db *gorm.DB,
	conf *config.Config,
) error {
	if !conf.Auth.Enabled {
		return nil
	}

	h.Register(&loginHandler{
		service: host.Service(),
	})
	h.Register(&callbackHandler{
		service:            host.Service(),
		successRedirectURL: conf.App.UIURL,
	})
	h.Register(&sessionHandler{service: host.Service()})

	tokenRepo := auth.NewRepo(db)
	h.Register(&listTokensHandler{repo: tokenRepo})
	h.Register(&createHandler{repo: tokenRepo})
	h.Register(&deleteHandler{repo: tokenRepo})
	return nil
}
