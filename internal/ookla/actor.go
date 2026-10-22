package ookla

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/henrywhitaker3/windowframe/v2/log"
	"github.com/henrywhitaker3/windowframe/v2/uuid"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
	"gorm.io/gorm"
)

const (
	Actor = "speedtest"
)

func RegisterActor(host *local.Host, db *gorm.DB) error {
	if err := host.RegisterActor(Actor, newActor(db)); err != nil {
		return fmt.Errorf("register speedtest actor: %w", err)
	}
	return nil
}

func newActor(db *gorm.DB) actor.Factory {
	return func(actorID string, service *actor.Service) actor.Actor {
		return &Speedtest{
			client: actor.NewActorClient[speedtestState](Actor, actorID, service),
			db:     db,
		}
	}
}

type Speedtest struct {
	client actor.Client[speedtestState]
	db     *gorm.DB
}

type speedtestState struct {
	Onboarded bool
	Result    *Result
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

	result := &Result{
		ID: uuid.MustOrdered(),
	}
	if err := gorm.G[Result](s.db).Create(ctx, result); err != nil {
		return nil, fmt.Errorf("store initial result: %w", err)
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

	result, err := gorm.G[Result](s.db).Where("id = ?", state.Result.ID).First(ctx)
	if err != nil {
		return fmt.Errorf("get refreshed result: %w", err)
	}
	state.Result = &result

	ctx = log.WithLogAttr(ctx, "result", result)

	if result.Complete {
		slog.WarnContext(ctx, "result already completed", "result", result)
		return nil
	}

	slog.InfoContext(ctx, "running speedtest")
	if err := Run(ctx, 0, &result); err != nil {
		slog.WarnContext(ctx, "speedtest failed", "error", err)
		return s.testFailed(ctx, &result, err)
	}
	ctx = log.WithLogAttr(ctx, "result", result)
	slog.InfoContext(ctx, "speedtest finished running")

	if err := s.saveResult(ctx, &result); err != nil {
		return fmt.Errorf("store updated test result: %w", err)
	}

	s.client.Halt()
	return s.client.DeleteState(ctx)
}

func (s *Speedtest) testFailed(ctx context.Context, result *Result, err error) error {
	result.Complete = true
	result.Error = new(fmt.Sprintf("%v", err.Error()))

	return s.saveResult(ctx, result)
}

func (s *Speedtest) saveResult(ctx context.Context, result *Result) error {
	if _, err := gorm.G[Result](s.db).
		Where("id = ?", result.ID).
		Select("*").
		Updates(ctx, *result); err != nil {
		return fmt.Errorf("store updated result: %w", err)
	}
	return nil
}
