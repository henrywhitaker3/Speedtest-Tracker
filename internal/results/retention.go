package results

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/italypaleale/francis/builtin/cronjob"
	"github.com/italypaleale/francis/host/local"
)

type retentionRepo interface {
	Forget(context.Context, time.Duration) (int, error)
}

func RegisterActors(host *local.Host, repo retentionRepo, period time.Duration) error {
	job, err := cronjob.New(
		"retention",
		cronjob.WithInterval(time.Hour),
		cronjob.WithJitter(time.Minute*5),
		cronjob.WithJob(func(ctx context.Context) error {
			if period <= 0 {
				slog.DebugContext(ctx, "skipping retention job", "reason", "zero retention period")
				return nil
			}
			rows, err := repo.Forget(ctx, period)
			if err != nil {
				return fmt.Errorf("forget results: %w", err)
			}
			slog.InfoContext(ctx, "deleted expired results", "count", rows)
			return nil
		}),
	)
	if err != nil {
		return fmt.Errorf("init cronjob: %w", err)
	}
	return host.RegisterBuiltInActor(job)
}
