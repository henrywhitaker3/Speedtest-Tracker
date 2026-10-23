package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/results"
	"github.com/henrywhitaker3/windowframe/v2/log"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
)

const (
	Actor = "speedtest"
)

type ResultRepo interface {
	Create(context.Context) (*results.Result, error)
	Fresh(context.Context, *results.Result) error
	Update(context.Context, *results.Result) error
}

func RegisterActor(host *local.Host, db ResultRepo, schedule string, server int) error {
	if err := host.RegisterActor(Actor, newActor(db, server)); err != nil {
		return fmt.Errorf("register speedtest actor: %w", err)
	}
	job, err := newTestCronjob(host.Service(), schedule)
	if err != nil {
		return err
	}
	if err := host.RegisterBuiltInActor(job); err != nil {
		return fmt.Errorf("register speedtest cronjob: %w", err)
	}
	return nil
}

func newActor(db ResultRepo, server int) actor.Factory {
	return func(actorID string, service *actor.Service) actor.Actor {
		return &Speedtest{
			client: actor.NewActorClient[speedtestState](Actor, actorID, service),
			db:     db,
			server: server,
		}
	}
}

type Speedtest struct {
	client actor.Client[speedtestState]
	db     ResultRepo
	server int
}

type speedtestState struct {
	Onboarded bool
	Result    *results.Result
}

const (
	Onboard = "onboard"
	RunTest = "run"
)

func (s *Speedtest) Invoke(ctx context.Context, method string, data actor.Envelope) (any, error) {
	switch method {
	case Onboard:
		return s.invokeOnboard(ctx)
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func (s *Speedtest) invokeOnboard(ctx context.Context) (any, error) {
	state, err := s.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get state: %w", err)
	}

	if state.Onboarded {
		return nil, fmt.Errorf("speedtest already onboarded")
	}

	result, err := s.db.Create(ctx)
	if err != nil {
		return nil, fmt.Errorf("init result record: %w", err)
	}

	state.Onboarded = true
	state.Result = result

	if err := s.client.SetState(ctx, state, nil); err != nil {
		return nil, fmt.Errorf("store speedtest state: %w", err)
	}

	if _, _, err := s.client.Dispatch(ctx, RunTest, nil); err != nil {
		return nil, fmt.Errorf("fire test job: %w", err)
	}

	return result, nil
}

func (s *Speedtest) Job(ctx context.Context, method string, data actor.Envelope) error {
	switch method {
	case RunTest:
		return s.runTest(ctx)
	default:
		return fmt.Errorf("unknown job %s", method)
	}
}

func (s *Speedtest) runTest(ctx context.Context) error {
	state, err := s.client.GetState(ctx)
	if err != nil {
		return fmt.Errorf("get speedtest state: %w", err)
	}

	if err := s.db.Fresh(ctx, state.Result); err != nil {
		return fmt.Errorf("get refreshed result: %w", err)
	}
	ctx = log.WithLogAttr(ctx, "result", state.Result)

	if state.Result.Complete {
		slog.WarnContext(ctx, "result already completed")
		return nil
	}

	slog.InfoContext(ctx, "running speedtest")
	if err := Run(ctx, s.server, state.Result); err != nil {
		slog.WarnContext(ctx, "speedtest failed", "error", err)
		return s.testFailed(ctx, state.Result, err)
	}
	ctx = log.WithLogAttr(ctx, "result", state.Result)
	slog.InfoContext(ctx, "speedtest finished running")

	if err := s.db.Update(ctx, state.Result); err != nil {
		return fmt.Errorf("store updated test result: %w", err)
	}

	s.client.Halt()
	return s.client.DeleteState(ctx)
}

func (s *Speedtest) testFailed(ctx context.Context, result *results.Result, err error) error {
	result.Complete = true
	result.Error = new(fmt.Sprintf("%v", err.Error()))

	return s.db.Update(ctx, result)
}
