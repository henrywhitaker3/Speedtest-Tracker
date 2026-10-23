// Package auth
package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/italypaleale/francis/actor"
	"github.com/italypaleale/francis/host/local"
)

const (
	Actor = "auth"
)

var errNoAuthenticatedSession = errors.New("no authenticated session")

type Auth struct {
	client      authStateClient
	config      config.Auth
	callbackURL string
	authorizer  authorizationURLBuilder
	refresher   tokenRefresher
	exchanger   tokenExchanger
}

const (
	refreshAhead        = time.Minute
	loginTransactionTTL = 10 * time.Minute
	sessionTTL          = time.Hour * 24 * 7
)

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
	OIDCResponse     OIDCResponse
	AuthenticatedAt  time.Time
	LoginTransaction *LoginTransaction
}

// LoginTransaction is persisted only until the OIDC callback is completed.
// Its verifier and nonce must never be sent to the browser.
type LoginTransaction struct {
	State     string
	Verifier  string
	Nonce     string
	ExpiresAt time.Time
}

type LoginResponse struct {
	RedirectURL string
}

type CompleteLoginRequest struct {
	Code  string
	State string
}

func RegisterActor(conf config.Auth, callbackURL string, host *local.Host) error {
	if err := host.RegisterActor(Actor, newAuthActor(conf, callbackURL)); err != nil {
		return fmt.Errorf("register auth actor: %w", err)
	}
	return nil
}

func CallbackURL(appURL string) string {
	return strings.TrimRight(appURL, "/") + "/v1/auth/callback"
}

const (
	Authenticate  = "authenticate"
	Session       = "session"
	BeginLogin    = "begin_login"
	CompleteLogin = "complete_login"
)

func (a *Auth) Invoke(ctx context.Context, method string, data actor.Envelope) (any, error) {
	switch method {
	case Authenticate:
		return a.invokeAuthenticate(ctx, data)
	case Session:
		return a.invokeSession(ctx)
	case BeginLogin:
		return a.invokeBeginLogin(ctx)
	case CompleteLogin:
		return a.invokeCompleteLogin(ctx, data)
	default:
		return nil, fmt.Errorf("unknown method %s", method)
	}
}

func (a *Auth) invokeBeginLogin(ctx context.Context) (*LoginResponse, error) {
	state, err := a.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth state: %w", err)
	}

	transaction, err := newLoginTransaction()
	if err != nil {
		return nil, fmt.Errorf("create login transaction: %w", err)
	}
	redirectURL, err := a.authorizer.AuthorizationURL(ctx, *transaction)
	if err != nil {
		return nil, fmt.Errorf("build OIDC authorization URL: %w", err)
	}
	state.LoginTransaction = transaction
	if err := a.client.SetState(
		ctx,
		state,
		&actor.SetStateOpts{TTL: loginTransactionTTL},
	); err != nil {
		return nil, fmt.Errorf("store login transaction: %w", err)
	}
	return &LoginResponse{RedirectURL: redirectURL}, nil
}

func newLoginTransaction() (*LoginTransaction, error) {
	state, err := randomValue(32)
	if err != nil {
		return nil, err
	}
	verifier, err := randomValue(32)
	if err != nil {
		return nil, err
	}
	nonce, err := randomValue(32)
	if err != nil {
		return nil, err
	}
	return &LoginTransaction{
		State:     state,
		Verifier:  verifier,
		Nonce:     nonce,
		ExpiresAt: time.Now().Add(loginTransactionTTL),
	}, nil
}

func randomValue(bytes int) (string, error) {
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func (a *Auth) invokeCompleteLogin(
	ctx context.Context,
	data actor.Envelope,
) (*AuthenticateResult, error) {
	if data == nil {
		return nil, errors.New("decode login callback: response is required")
	}
	var request CompleteLoginRequest
	if err := data.Decode(&request); err != nil {
		return nil, fmt.Errorf("decode login callback: %w", err)
	}
	if request.Code == "" || request.State == "" {
		return nil, errors.New("decode login callback: code and state are required")
	}

	state, err := a.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth state: %w", err)
	}
	transaction := state.LoginTransaction
	if transaction == nil || !transaction.ExpiresAt.After(time.Now()) {
		return nil, errors.New("login transaction has expired")
	}
	if subtle.ConstantTimeCompare([]byte(request.State), []byte(transaction.State)) != 1 {
		return nil, errors.New("invalid login state")
	}

	response, err := a.exchanger.Exchange(
		ctx,
		request.Code,
		transaction.Verifier,
		transaction.Nonce,
	)
	if err != nil {
		return nil, fmt.Errorf("exchange OIDC authorization code: %w", err)
	}
	if response.Subject == "" || !response.ExpiresAt.After(time.Now()) {
		return nil, errors.New("OIDC response is missing a valid subject or expiry")
	}

	state.OIDCResponse = response
	state.AuthenticatedAt = time.Now().UTC()
	state.LoginTransaction = nil
	if err := a.client.SetState(ctx, state, &actor.SetStateOpts{
		TTL: sessionTTL,
	}); err != nil {
		return nil, fmt.Errorf("store authenticated session: %w", err)
	}
	result := sessionResult(response)
	return &result, nil
}

func (a *Auth) invokeSession(ctx context.Context) (any, error) {
	state, err := a.client.GetState(ctx)
	if err != nil {
		return nil, fmt.Errorf("get auth state: %w", err)
	}
	if state.OIDCResponse.Subject == "" {
		return nil, errNoAuthenticatedSession
	}

	now := time.Now()
	if state.OIDCResponse.ExpiresAt.Before(now.Add(refreshAhead)) &&
		state.OIDCResponse.RefreshToken != "" {
		refreshed, err := a.refresher.Refresh(ctx, state.OIDCResponse)
		if err == nil {
			state.OIDCResponse = refreshed
			if err := a.client.SetState(ctx, state, &actor.SetStateOpts{
				TTL: sessionTTL,
			}); err != nil {
				return nil, fmt.Errorf("store refreshed auth state: %w", err)
			}
		} else if !state.OIDCResponse.ExpiresAt.After(now) {
			a.client.Halt()
			_ = a.client.DeleteState(ctx)
			return nil, fmt.Errorf("refresh expired session: %w", err)
		}
	}

	if !state.OIDCResponse.ExpiresAt.After(now) {
		a.client.Halt()
		_ = a.client.DeleteState(ctx)
		return nil, fmt.Errorf("expired session")
	}

	return sessionResult(state.OIDCResponse), nil
}

func sessionResult(response OIDCResponse) AuthenticateResult {
	return AuthenticateResult{
		Name:      response.Name,
		Subject:   response.Subject,
		Email:     response.Email,
		ExpiresAt: response.ExpiresAt,
		Groups:    response.Groups,
	}
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
	Subject   string    `json:"subject"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Groups    []string  `json:"groups"`
	ExpiresAt time.Time `json:"expires_at"`
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
		TTL: sessionTTL,
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

func newAuthActor(conf config.Auth, callbackURL string) actor.Factory {
	return func(actorID string, service *actor.Service) actor.Actor {
		return &Auth{
			client:      actor.NewActorClient[authState](Actor, actorID, service),
			config:      conf,
			callbackURL: callbackURL,
			authorizer: oidcAuthorizationURLBuilder{
				config:      conf,
				callbackURL: callbackURL,
			},
			exchanger: oidcTokenExchanger{
				config:      conf,
				callbackURL: callbackURL,
			},
			refresher: oidcTokenRefresher{config: conf},
		}
	}
}
