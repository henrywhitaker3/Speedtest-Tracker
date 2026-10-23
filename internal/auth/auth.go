// Package auth
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
)

const (
	Actor = "auth"
)

type Auth struct {
	client authStateClient
	config config.Auth
}

// authStateClient is the small part of an actor client that the auth actor
// needs. Keeping this narrow also makes the persistence behaviour testable
// without starting a runtime host.
type authStateClient interface {
	GetState(context.Context) (authState, error)
	SetState(context.Context, authState, *actor.SetStateOpts) error
	Halt()
	DeleteState(context.Context) error
}

// authState is stored per authenticated browser session. OIDCResponse contains
// credentials, so it must never be returned from an HTTP handler or actor
// invocation response.
type authState struct {
	OIDCResponse    OIDCResponse
	AuthenticatedAt time.Time
}

func RegisterActor(conf config.Auth, host *local.Host) error {
	if err := host.RegisterActor(Actor, newAuthActor(conf)); err != nil {
		return fmt.Errorf("register auth actor: %w", err)
	}
	return nil
}

const (
	Authenticate = "authenticate"
	Session      = "session"
)

func (a *Auth) Invoke(ctx context.Context, method string, data actor.Envelope) (any, error) {
	switch method {
	case Authenticate:
		return a.invokeAuthenticate(ctx, data)
	case Session:
		return a.invokeSession(ctx)
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func (a *Auth) invokeSession(ctx context.Context) (any, error) {
	state, err := a.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth state: %w", err)
	}
	if state.OIDCResponse.ExpiresAt.Before(time.Now()) {
		a.client.Halt()
		_ = a.client.DeleteState(ctx)
		return nil, fmt.Errorf("expired session")
	}

	return AuthenticateResult{
		Name:      state.OIDCResponse.Name,
		Subject:   state.OIDCResponse.Subject,
		Email:     state.OIDCResponse.Email,
		ExpiresAt: state.OIDCResponse.ExpiresAt,
		Groups:    state.OIDCResponse.Groups,
	}, nil
}

// OIDCResponse is the already-validated response received from an OIDC issuer.
// The HTTP callback is responsible for exchanging the authorization code and
// verifying the ID token before it invokes this actor.
type OIDCResponse struct {
	AccessToken  string
	TokenType    string
	RefreshToken string
	IDToken      string
	Scope        string
	ExpiresAt    time.Time

	Subject string
	Email   string
	Name    string
	Groups  []string
}

// AuthenticateResult is safe to return to callers. In particular, it does not
// expose the issuer's access, refresh, or ID tokens.
type AuthenticateResult struct {
	Subject   string
	Email     string
	Name      string
	Groups    []string
	ExpiresAt time.Time
}

func (a *Auth) invokeAuthenticate(
	ctx context.Context,
	data actor.Envelope,
) (*AuthenticateResult, error) {
	if data == nil {
		return nil, fmt.Errorf("decode auth response: response is required")
	}

	var response OIDCResponse
	if err := data.Decode(&response); err != nil {
		return nil, fmt.Errorf("decode auth response: %w", err)
	}
	if response.Subject == "" {
		return nil, fmt.Errorf("decode auth response: subject is required")
	}

	state, err := a.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth state: %w", err)
	}

	state.OIDCResponse = response
	state.AuthenticatedAt = time.Now().UTC()
	if err := a.client.SetState(ctx, state, &actor.SetStateOpts{
		TTL: time.Until(response.ExpiresAt),
	}); err != nil {
		return nil, fmt.Errorf("store auth state: %w", err)
	}

	return &AuthenticateResult{
		Subject:   response.Subject,
		Email:     response.Email,
		Name:      response.Name,
		Groups:    response.Groups,
		ExpiresAt: response.ExpiresAt,
	}, nil
}

func newAuthActor(conf config.Auth) actor.Factory {
	return func(actorID string, service *actor.Service) actor.Actor {
		return &Auth{
			client: actor.NewActorClient[authState](Actor, actorID, service),
			config: conf,
		}
	}
}
