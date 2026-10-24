package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/notifications"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/log"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
)

type workflowRepo interface {
	Fresh(context.Context, *results.Result) error
	Update(context.Context, *results.Result) error
}

type runFunc func(context.Context, int, *results.Result) error

type notifier interface {
	ShouldNotify() bool
	Notify(context.Context, notifications.Phase, results.Result) error
}

type WorkflowOpts struct {
	Host     *local.Host
	Repo     workflowRepo
	Runner   runFunc
	Notifier notifier
}

func buildWorkflow(opts WorkflowOpts) (*workflow.Workflow, error) {
	return workflow.New(
		"speedtest",
		workflow.WithVersion(3),
		workflow.WithSteps(
			workflow.Step(
				"notify",
				workflow.WithRun(shouldNotify(opts.Notifier)),
			),
			workflow.Step(
				"pre-run notifications",
				workflow.WithRun(preRunNotify(opts.Notifier)),
				workflow.WithSkipIf("notify", false),
				workflow.WithOptional(),
			),
			workflow.Step(
				"run test",
				workflow.WithRun(workflowRun(opts.Repo, opts.Runner)),
			),
			workflow.Step(
				"store result",
				workflow.WithRun(storeResult(opts.Repo)),
			),
			workflow.Step(
				"post-run notifications",
				workflow.WithRun(postRunNotify(opts.Notifier)),
				workflow.WithOptional(),
				workflow.WithSkipIf("notify", false),
			),
		),
	)
}

type RunWorkflow struct {
	ServerID int
	Result   results.Result
}

func workflowRun(repo workflowRepo, runner runFunc) workflow.RunFunc {
	return func(ctx context.Context, task workflow.Task) (any, error) {
		var opts RunWorkflow
		if err := task.DecodeInput(&opts); err != nil {
			return nil, fmt.Errorf("decode result model: %w", err)
		}

		result := opts.Result
		if err := repo.Fresh(ctx, &result); err != nil {
			return nil, fmt.Errorf("refresh result: %w", err)
		}
		ctx = log.WithLogAttr(ctx, "result", result.ID)
		if result.Complete {
			slog.WarnContext(ctx, "speedest already complete")
			return result, nil
		}

		err := runner(ctx, opts.ServerID, &result)
		result.Complete = true
		if err != nil {
			slog.WarnContext(ctx, "speedtest failed", "error", err)
			result.Error = new(err.Error())
		}
		slog.InfoContext(ctx, "speedtest finished running")

		return result, nil
	}
}

func storeResult(repo workflowRepo) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		var result results.Result
		if err := t.DecodeOutput("run test", &result); err != nil {
			return nil, fmt.Errorf("decode run output: %w", err)
		}

		if err := repo.Update(ctx, &result); err != nil {
			return nil, fmt.Errorf("store workflow result: %w", err)
		}

		return result, nil
	}
}

func shouldNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		return n.ShouldNotify(), nil
	}
}

func preRunNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		var result RunWorkflow
		if err := t.DecodeInput(&result); err != nil {
			return nil, fmt.Errorf("decode input: %w", err)
		}
		return nil, n.Notify(ctx, notifications.PreRun, result.Result)
	}
}

func postRunNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (output any, err error) {
		var result results.Result
		if err := t.DecodeOutput("store result", &result); err != nil {
			return nil, fmt.Errorf("decode store result output: %w", err)
		}
		return nil, n.Notify(ctx, notifications.PostRun, result)
	}
}
