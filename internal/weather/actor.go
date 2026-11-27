package weather

import (
	"fmt"

	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
)

func RegisterActor(host *local.Host, p HistoricalProvider, repo historicalWeatherRepo) error {
	if p != nil {
		if err := host.RegisterSingletonActor(
			Actor,
			func(actorID string, service *actor.Service) actor.Actor {
				return &historicalActor{
					client:   actor.NewActorClient[historicalActorState](Actor, actorID, service),
					provider: p,
					repo:     repo,
				}
			},
		); err != nil {
			return fmt.Errorf("register historical actor: %w", err)
		}
	}

	return nil
}
