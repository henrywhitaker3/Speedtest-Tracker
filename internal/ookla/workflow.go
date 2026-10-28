package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/notifications"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/log"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
	"go.opentelemetry.io/otel"
)

type workflowRepo interface {
	Fresh(context.Context, *results.Result) error
	Update(context.Context, *results.Result) error
}

type testRunner interface {
	Run(context.Context, int, *results.Result) error
	PickServer(context.Context) (int, error)
}

type notifier interface {
	Invoke(
		context.Context,
		string,
		string,
		string,
		any,
		...actor.InvokeOption,
	) (actor.Envelope, error)
}

type WorkflowOpts struct {
	Host     *local.Host
	Repo     workflowRepo
	Runner   testRunner
	Notifier notifier
}

func buildWorkflow(opts WorkflowOpts) (*workflow.Workflow, error) {
	return workflow.New(
		"speedtest",
		workflow.WithVersion(7),
		workflow.WithUnknownVersionPolicy(workflow.FailUnknownVersion),
		workflow.WithMeter(otel.Meter("speedtest")),
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
				"pick server",
				workflow.WithRun(pickServer(opts.Runner, opts.Repo)),
			),
			workflow.Step(
				"run test",
				workflow.WithRun(workflowRun(opts.Repo, opts.Runner)),
				workflow.WithInputFrom("pick server"),
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
	Result results.Result
}

func pickServer(runner testRunner, repo workflowRepo) workflow.RunFunc {
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
			return 0, nil
		}

		server, err := runner.PickServer(ctx)
		if err != nil {
			result.Complete = true
			result.Error = new(err.Error())
			_ = repo.Update(ctx, &result)
		}

		return server, err
	}
}

func workflowRun(repo workflowRepo, runner testRunner) workflow.RunFunc {
	return func(ctx context.Context, task workflow.Task) (any, error) {
		var opts RunWorkflow
		if err := task.DecodeInput(&opts); err != nil {
			return nil, fmt.Errorf("decode result model: %w", err)
		}

		var server int
		if err := task.DecodeOutput("pick server", &server); err != nil {
			return nil, fmt.Errorf("get pick server output: %w", err)
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

		err := runner.Run(ctx, server, &result)
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

// c.Request().Context(),
// authactor.Actor,
// sessionID,
// authactor.BeginLogin,
// nil,

func shouldNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		data, err := n.Invoke(
			ctx,
			notifications.Actor,
			uuid.MustNew().String(),
			notifications.ShouldNotify,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("invoke notifications actor: %w", err)
		}
		var should bool
		_ = data.Decode(&should)
		return should, nil
	}
}

func preRunNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		var result RunWorkflow
		if err := t.DecodeInput(&result); err != nil {
			return nil, fmt.Errorf("decode input: %w", err)
		}

		data, err := n.Invoke(
			ctx,
			notifications.Actor,
			result.Result.ID.String(),
			notifications.Notify,
			notifications.NotifyInput{
				Result: result.Result,
				Phase:  notifications.PreRun,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("invoke notifications actor: %w", err)
		}
		var errors *notifications.NotificationError
		if err := data.Decode(&errors); err != nil {
			return nil, fmt.Errorf("decode actor response: %w", err)
		}
		if errors == nil || errors.Nil() {
			return nil, nil
		}

		return nil, errors
	}
}

func postRunNotify(n notifier) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (output any, err error) {
		var result results.Result
		if err := t.DecodeOutput("store result", &result); err != nil {
			return nil, fmt.Errorf("decode input: %w", err)
		}

		data, err := n.Invoke(
			ctx,
			notifications.Actor,
			result.ID.String(),
			notifications.Notify,
			notifications.NotifyInput{
				Result: result,
				Phase:  notifications.PostRun,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("invoke notifications actor: %w", err)
		}
		var errors *notifications.NotificationError
		if err := data.Decode(&errors); err != nil {
			return nil, fmt.Errorf("decode actor response: %w", err)
		}
		if errors == nil || errors.Nil() {
			return nil, nil
		}

		return nil, errors
	}
}
