// Package speedtest
package speedtest

import (
	"fmt"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/http"
	"github.com/italypaleale/francis/builtin/ratelimit"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
	"gorm.io/gorm"
)

func Register(
	h *http.HTTP,
	db *gorm.DB,
	host *local.Host,
	runner *workflow.WorkflowService,
	serverID int,
) error {
	runLimiter, err := ratelimit.New(
		"run",
		ratelimit.WithPer(time.Minute),
		ratelimit.WithRate(3),
	)
	if err != nil {
		return fmt.Errorf("init run limiter: %w", err)
	}
	if err := host.RegisterBuiltInActor(runLimiter); err != nil {
		return fmt.Errorf("register run limiter: %w", err)
	}

	repo := results.NewRepo(db)
	h.Register(&runHandler{
		repo:     repo,
		workflow: runner,
		limiter:  runLimiter.Service(host.Service()),
	})
	h.Register(&getHandler{
		repo:   repo,
		status: runner,
	})
	h.Register(&listHandler{repo: repo})
	h.Register(&latestHandler{repo: repo})
	h.Register(&timeseriesHandler{repo: repo})
	h.Register(&purgeHandler{repo: repo})
	h.Register(&deleteHandler{repo: repo})

	return nil
}
