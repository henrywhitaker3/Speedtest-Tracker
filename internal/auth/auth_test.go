package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/italypaleale/francis/actor"
)

func TestInvokeAuthenticateStoresOIDCResponse(t *testing.T) {
	expiresAt := time.Date(2026, 10, 23, 13, 0, 0, 0, time.UTC)
	client := &fakeStateClient{}
	auth := &Auth{client: client}

	result, err := auth.invokeAuthenticate(context.Background(), envelope{value: OIDCResponse{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		Scope:        "openid profile email",
		ExpiresAt:    expiresAt,
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		Name:         "Person Example",
		Groups:       []string{"speedtest-admins"},
	}})
	require.NoError(t, err)

	assert.Equal(t, OIDCResponse{
		AccessToken:  "access-token",
		TokenType:    "Bearer",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		Scope:        "openid profile email",
		ExpiresAt:    expiresAt,
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		Name:         "Person Example",
		Groups:       []string{"speedtest-admins"},
	}, client.state.OIDCResponse)

	assert.False(t, client.state.AuthenticatedAt.IsZero())
	require.NotNil(t, client.setOpts)
	assert.Equal(t, sessionTTL, client.setOpts.TTL)

	want := &AuthenticateResult{
		Subject:   "issuer-user-id",
		Email:     "person@example.com",
		Name:      "Person Example",
		Groups:    []string{"speedtest-admins"},
		ExpiresAt: expiresAt,
	}
	assert.Equal(t, want, result)
}

func TestNewAuthActorConfiguresOIDCDependencies(t *testing.T) {
	conf := config.Auth{IssuerURL: "https://issuer.example", ClientID: "client-id"}
	callbackURL := "https://speedtest.example/v1/auth/callback"

	auth, ok := newAuthActor(conf, callbackURL, nil)("session-id", nil).(*Auth)
	require.True(t, ok)

	_, ok = auth.authorizer.(oidcAuthorizationURLBuilder)
	require.True(t, ok)

	_, ok = auth.exchanger.(oidcTokenExchanger)
	require.True(t, ok)
	_, ok = auth.refresher.(oidcTokenRefresher)
	require.True(t, ok)
}

func TestInvokeAuthenticateDoesNotStoreInvalidResponse(t *testing.T) {
	client := &fakeStateClient{}
	auth := &Auth{client: client}

	_, err := auth.invokeAuthenticate(context.Background(), envelope{value: OIDCResponse{}})
	require.Error(t, err)
	assert.False(t, client.setCalled)
}

func TestInvokeAuthenticateReturnsStoreError(t *testing.T) {
	storeErr := errors.New("store unavailable")
	client := &fakeStateClient{setErr: storeErr}
	auth := &Auth{client: client}

	_, err := auth.invokeAuthenticate(
		context.Background(),
		envelope{value: OIDCResponse{Subject: "issuer-user-id"}},
	)
	assert.ErrorIs(t, err, storeErr)
}

func TestInvokeBeginLoginStoresTransaction(t *testing.T) {
	client := &fakeStateClient{}
	authorizer := &authorizationURLBuilderStub{url: "https://issuer.example/authorize"}
	auth := &Auth{client: client, authorizer: authorizer}

	response, err := auth.invokeBeginLogin(context.Background())
	require.NoError(t, err)
	assert.Equal(t, authorizer.url, response.RedirectURL)
	transaction := client.state.LoginTransaction
	require.NotNil(t, transaction)
	assert.NotEmpty(t, transaction.State)
	assert.NotEmpty(t, transaction.Verifier)
	assert.NotEmpty(t, transaction.Nonce)
	assert.True(t, transaction.ExpiresAt.After(time.Now()))
	assert.True(t, authorizer.called)
	assert.True(t, client.setCalled)
	require.NotNil(t, client.setOpts)
	assert.Equal(t, loginTransactionTTL, client.setOpts.TTL)
}

func TestInvokeCompleteLoginExchangesAndStoresSession(t *testing.T) {
	transaction := &LoginTransaction{
		State:     "state",
		Verifier:  "verifier",
		Nonce:     "nonce",
		ExpiresAt: time.Now().Add(loginTransactionTTL),
	}
	expiresAt := time.Now().Add(time.Hour).UTC()
	response := OIDCResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		ExpiresAt:    expiresAt,
	}
	client := &fakeStateClient{state: authState{LoginTransaction: transaction}}
	exchanger := &tokenExchangerStub{response: response}
	auth := &Auth{client: client, exchanger: exchanger}

	result, err := auth.invokeCompleteLogin(
		context.Background(),
		envelope{value: CompleteLoginRequest{
			Code: "authorization-code", State: "state",
		}},
	)
	require.NoError(t, err)
	assert.True(t, exchanger.called)
	assert.Equal(t, "authorization-code", exchanger.code)
	assert.Equal(t, "verifier", exchanger.verifier)
	assert.Equal(t, "nonce", exchanger.nonce)
	assert.Nilf(
		t,
		client.state.LoginTransaction,
		"login transaction was not cleared after successful authentication",
	)
	assert.Equal(t, response, client.state.OIDCResponse)
	require.NotNil(t, client.setOpts)
	assert.Equal(t, sessionTTL, client.setOpts.TTL)
	assert.Equal(t, response.Subject, result.Subject)
	assert.Equal(t, response.ExpiresAt, result.ExpiresAt)
}

func TestInvokeCompleteLoginRejectsMismatchedState(t *testing.T) {
	client := &fakeStateClient{state: authState{LoginTransaction: &LoginTransaction{
		State: "expected-state", ExpiresAt: time.Now().Add(loginTransactionTTL),
	}}}
	exchanger := &tokenExchangerStub{}
	auth := &Auth{client: client, exchanger: exchanger}

	_, err := auth.invokeCompleteLogin(context.Background(), envelope{value: CompleteLoginRequest{
		Code: "authorization-code", State: "wrong-state",
	}})
	require.Error(t, err)
	assert.False(t, exchanger.called)
	assert.False(t, client.setCalled)
}

func TestInvokeSessionReturnsSafeSessionDetails(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour).UTC()
	client := &fakeStateClient{state: authState{OIDCResponse: OIDCResponse{
		AccessToken:  "access-token",
		RefreshToken: "refresh-token",
		IDToken:      "id-token",
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		Name:         "Person Example",
		Groups:       []string{"speedtest-admins"},
		ExpiresAt:    expiresAt,
	}}}
	auth := &Auth{client: client}

	result, err := auth.invokeSession(context.Background())
	require.NoError(t, err)

	want := AuthenticateResult{
		Subject:   "issuer-user-id",
		Email:     "person@example.com",
		Name:      "Person Example",
		Groups:    []string{"speedtest-admins"},
		ExpiresAt: expiresAt,
	}
	assert.Equal(t, want, result)
	assert.False(t, client.haltCalled)
	assert.False(t, client.deleteCalled)
}

func TestInvokeSessionRejectsStateWithoutAuthenticatedSession(t *testing.T) {
	client := &fakeStateClient{}
	auth := &Auth{client: client}

	_, err := auth.invokeSession(context.Background())
	assert.ErrorIs(t, err, errNoAuthenticatedSession)
	assert.False(t, client.haltCalled)
	assert.False(t, client.deleteCalled)
}

func TestInvokeSessionRefreshesExpiringSession(t *testing.T) {
	oldExpiry := time.Now().Add(30 * time.Second).UTC()
	newExpiry := time.Now().Add(time.Hour).UTC()
	client := &fakeStateClient{state: authState{OIDCResponse: OIDCResponse{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		ExpiresAt:    oldExpiry,
	}}}
	refresher := &tokenRefresherStub{response: OIDCResponse{
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		Subject:      "issuer-user-id",
		Email:        "person@example.com",
		ExpiresAt:    newExpiry,
	}}
	auth := &Auth{client: client, refresher: refresher}

	result, err := auth.invokeSession(context.Background())
	require.NoError(t, err)
	assert.True(t, refresher.called)
	assert.Equal(t, "new-access-token", client.state.OIDCResponse.AccessToken)
	require.True(t, client.setCalled)
	require.NotNil(t, client.setOpts)
	assert.Equal(t, sessionTTL, client.setOpts.TTL)
	got := result.(AuthenticateResult)
	assert.Equal(t, newExpiry, got.ExpiresAt)
}

func TestInvokeSessionInvalidatesExpiredSessionWhenRefreshFails(t *testing.T) {
	refreshErr := errors.New("issuer unavailable")
	client := &fakeStateClient{state: authState{OIDCResponse: OIDCResponse{
		RefreshToken: "refresh-token",
		Subject:      "issuer-user-id",
		ExpiresAt:    time.Now().Add(-time.Hour),
	}}}
	auth := &Auth{client: client, refresher: &tokenRefresherStub{err: refreshErr}}

	_, err := auth.invokeSession(context.Background())
	assert.ErrorIs(t, err, refreshErr)
	assert.True(t, client.haltCalled)
	assert.True(t, client.deleteCalled)
}

func TestInvokeSessionExpiresAndDeletesState(t *testing.T) {
	client := &fakeStateClient{state: authState{OIDCResponse: OIDCResponse{
		Subject:   "issuer-user-id",
		ExpiresAt: time.Now().Add(-time.Hour),
	}}}
	auth := &Auth{client: client}

	_, err := auth.invokeSession(context.Background())
	require.Error(t, err)
	assert.Equal(t, "expired session", err.Error())
	assert.True(t, client.haltCalled)
	assert.True(t, client.deleteCalled)
}

func TestInvokeSessionReturnsStateReadError(t *testing.T) {
	readErr := errors.New("state unavailable")
	auth := &Auth{client: &fakeStateClient{getErr: readErr}}

	_, err := auth.invokeSession(context.Background())
	assert.ErrorIs(t, err, readErr)
}

type fakeStateClient struct {
	state        authState
	getErr       error
	setErr       error
	deleteErr    error
	setCalled    bool
	setOpts      *actor.SetStateOpts
	haltCalled   bool
	deleteCalled bool
}

func (f *fakeStateClient) GetState(context.Context) (authState, error) {
	return f.state, f.getErr
}

func (f *fakeStateClient) SetState(
	_ context.Context,
	state authState,
	opts *actor.SetStateOpts,
) error {
	f.setCalled = true
	f.setOpts = opts
	if f.setErr != nil {
		return f.setErr
	}
	f.state = state
	return nil
}

type tokenRefresherStub struct {
	response OIDCResponse
	err      error
	called   bool
}

type tokenExchangerStub struct {
	response OIDCResponse
	err      error
	code     string
	verifier string
	nonce    string
	called   bool
}

type authorizationURLBuilderStub struct {
	url    string
	err    error
	called bool
}

func (a *authorizationURLBuilderStub) AuthorizationURL(
	_ context.Context,
	_ LoginTransaction,
) (string, error) {
	a.called = true
	return a.url, a.err
}

func (e *tokenExchangerStub) Exchange(
	_ context.Context,
	code, verifier, nonce string,
) (OIDCResponse, error) {
	e.called = true
	e.code = code
	e.verifier = verifier
	e.nonce = nonce
	return e.response, e.err
}

func (r *tokenRefresherStub) Refresh(_ context.Context, _ OIDCResponse) (OIDCResponse, error) {
	r.called = true
	return r.response, r.err
}

func (f *fakeStateClient) Halt() {
	f.haltCalled = true
}

func (f *fakeStateClient) DeleteState(context.Context) error {
	f.deleteCalled = true
	return f.deleteErr
}

type envelope struct {
	value any
	err   error
}

func (e envelope) Decode(into any) error {
	if e.err != nil {
		return e.err
	}
	switch target := into.(type) {
	case *OIDCResponse:
		*target = e.value.(OIDCResponse)
	case *CompleteLoginRequest:
		*target = e.value.(CompleteLoginRequest)
	default:
		return errors.New("unexpected decode target")
	}
	return nil
}
