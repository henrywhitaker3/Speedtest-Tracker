package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/cronjob"
)

func newTestCronjob(service *actor.Service, schedule string) (*cronjob.CronJob, error) {
	job, err := cronjob.New(
		Actor,
		cronjob.WithCron(schedule),
		cronjob.WithJob(runScheduleTest(service)),
	)
	if err != nil {
		return nil, fmt.Errorf("build speedtest cronjob: %w", err)
	}
	return job, nil
}

func runScheduleTest(service *actor.Service) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		slog.InfoContext(ctx, "triggering scheduled speedtest")
		id, _ := uuid.New()
		if _, err := service.Invoke(ctx, Actor, id.String(), Onboard, nil); err != nil {
			return fmt.Errorf("onboard new test: %w", err)
		}

		return nil
	}
}
