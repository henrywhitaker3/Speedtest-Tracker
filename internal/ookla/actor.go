package ookla

import (
	"context"
	"fmt"
	"time"

	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/builtin/workflow"
	"github.com/italypaleale/francis/host/local"
)

type RegisterActorOpts struct {
	Host         *local.Host
	CronjobRepo  resultCreator
	WorkflowRepo workflowRepo
	Schedule     string
	Runner       testRunner
	Lister       lister
	Notifier     notifier
}

func RegisterActor(opts RegisterActorOpts) (*workflow.WorkflowService, error) {
	workflow, err := buildWorkflow(WorkflowOpts{
		Host:     opts.Host,
		Runner:   opts.Runner,
		Repo:     opts.WorkflowRepo,
		Notifier: opts.Notifier,
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
	)
	if err != nil {
		return nil, err
	}
	if err := opts.Host.RegisterBuiltInActor(job); err != nil {
		return nil, fmt.Errorf("register speedtest cronjob: %w", err)
	}

	if err := opts.Host.RegisterSingletonActor(
		Actor,
		func(actorID string, service *actor.Service) actor.Actor {
			return &serverListerActor{
				client: actor.NewActorClient[serverListerState](Actor, actorID, service),
				lister: opts.Lister,
			}
		},
	); err != nil {
		return nil, fmt.Errorf("register server lister: %w", err)
	}

	return workflow.Service(opts.Host.Service()), nil
}

type serverListerState struct {
	Envelope *serverEnvelope
}

type lister interface {
	listServers(ctx context.Context) (serverEnvelope, error)
}

type serverListerActor struct {
	client actor.Client[serverListerState]
	lister lister
}

const (
	Actor       = "server-lister"
	ListServers = "list-servers"
)

func (s *serverListerActor) Invoke(
	ctx context.Context,
	method string,
	_ actor.Envelope,
) (any, error) {
	switch method {
	case ListServers:
		return s.invokeListServers(ctx)
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func (s *serverListerActor) invokeListServers(ctx context.Context) ([]Server, error) {
	state, err := s.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get lister state: %w", err)
	}

	if state.Envelope != nil {
		if !state.Envelope.Timestamp.Before(time.Now().Add(-time.Minute * 30)) {
			return state.Envelope.Servers, nil
		}
	}

	envelope, err := s.lister.listServers(ctx)
	if err != nil {
		return nil, err
	}

	state.Envelope = &envelope
	if err := s.client.SetState(ctx, state, nil); err != nil {
		return nil, fmt.Errorf("store lister state: %w", err)
	}

	return state.Envelope.Servers, nil
}
