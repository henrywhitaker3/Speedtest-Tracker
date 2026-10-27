package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/builtin/cronjob"
	"github.com/italypaleale/francis/builtin/workflow"
)

type resultCreator interface {
	Create(context.Context) (*results.Result, error)
	Delete(context.Context, uuid.UUID) error
}

type worflowRunner interface {
	Start(context.Context, any, ...workflow.StartOption) (string, bool, error)
}

func newTestCronjob(
	repo resultCreator,
	runner worflowRunner,
	schedule string,
) (*cronjob.CronJob, error) {
	job, err := cronjob.New(
		"speedtest",
		cronjob.WithCron(schedule),
		cronjob.WithJob(runScheduleTest(repo, runner)),
	)
	if err != nil {
		return nil, fmt.Errorf("build speedtest cronjob: %w", err)
	}
	return job, nil
}

func runScheduleTest(
	repo resultCreator,
	runner worflowRunner,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		slog.InfoContext(ctx, "triggering scheduled speedtest")

		result, err := repo.Create(ctx)
		if err != nil {
			return fmt.Errorf("create initial result record: %w", err)
		}

		if _, _, err := runner.Start(
			ctx,
			RunWorkflow{
				Result: *result,
			},
			workflow.WithInstanceID(result.ID.String()),
		); err != nil {
			_ = repo.Delete(ctx, result.ID)
			return fmt.Errorf("start workflow: %w", err)
		}

		return nil
	}
}
