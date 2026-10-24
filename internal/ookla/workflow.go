package ookla

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

type WorkflowOpts struct {
	Host        *local.Host
	Repo        workflowRepo
	Runner      runFunc
	Healthcheck string
}

func buildWorkflow(opts WorkflowOpts) (*workflow.Workflow, error) {
	return workflow.New(
		"speedtest",
		workflow.WithVersion(2),
		workflow.WithSteps(
			workflow.Step(
				"healthcheck",
				workflow.WithRun(shouldRunHealthcheck(opts.Healthcheck)),
			),
			workflow.Step(
				"start healthcheck",
				workflow.WithRun(startHealthcheck(opts.Healthcheck)),
				workflow.WithSkipIf("healthcheck", false),
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
				"finish healthcheck",
				workflow.WithRun(finishHealthcheck(opts.Healthcheck)),
				workflow.WithOptional(),
				workflow.WithSkipIf("healthcheck", false),
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

func shouldRunHealthcheck(url string) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		return url != "", nil
	}
}

func startHealthcheck(url string) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (any, error) {
		return nil, sendHealthcheck(ctx, url, healthcheckStart)
	}
}

func finishHealthcheck(url string) workflow.RunFunc {
	return func(ctx context.Context, t workflow.Task) (output any, err error) {
		var result results.Result
		if err := t.DecodeOutput("store result", &result); err != nil {
			return nil, fmt.Errorf("decode store result output: %w", err)
		}

		action := healthcheckFail
		if result.Error == nil {
			action = healthcheckSuccess
		}

		return nil, sendHealthcheck(ctx, url, action)
	}
}

type healthcheckAction = string

const (
	healthcheckStart   healthcheckAction = "/start"
	healthcheckFail    healthcheckAction = "/fail"
	healthcheckSuccess healthcheckAction = ""
)

func sendHealthcheck(ctx context.Context, url string, action healthcheckAction) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s%s", url, action),
		nil,
	)
	if err != nil {
		return fmt.Errorf("build http request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do healthcheck request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d", resp.StatusCode)
	}

	return nil
}
