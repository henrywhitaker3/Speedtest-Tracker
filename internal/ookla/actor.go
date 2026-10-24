package ookla

import (
	"fmt"

	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
)

type RegisterActorOpts struct {
	Host         *local.Host
	CronjobRepo  resultCreator
	WorkflowRepo workflowRepo
	Schedule     string
	ServerID     int
	Runner       runFunc
	Healthcheck  string
}

func RegisterActor(opts RegisterActorOpts) (*workflow.WorkflowService, error) {
	workflow, err := buildWorkflow(WorkflowOpts{
		Host:        opts.Host,
		Runner:      opts.Runner,
		Repo:        opts.WorkflowRepo,
		Healthcheck: opts.Healthcheck,
	})
	if err != nil {
		return nil, fmt.Errorf("build workflow: %w", err)
	}
	if err := opts.Host.RegisterBuiltInActor(workflow); err != nil {
		return nil, fmt.Errorf("register workflow: %w", err)
	}

	job, err := newTestCronjob(
		opts.CronjobRepo,
		workflow.Service(opts.Host.Service()),
		opts.Schedule,
		opts.ServerID,
	)
	if err != nil {
		return nil, err
	}
	if err := opts.Host.RegisterBuiltInActor(job); err != nil {
		return nil, fmt.Errorf("register speedtest cronjob: %w", err)
	}

	return workflow.Service(opts.Host.Service()), nil
}
