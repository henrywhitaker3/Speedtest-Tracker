package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

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
	if err != nil {
		t.Fatalf("invokeAuthenticate() error = %v", err)
	}

	if !reflect.DeepEqual(client.state.OIDCResponse, OIDCResponse{
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
	}) {
		t.Errorf("stored OIDC response = %#v", client.state.OIDCResponse)
	}
	if client.state.AuthenticatedAt.IsZero() {
		t.Error("AuthenticatedAt was not set")
	}

	want := &AuthenticateResult{
		Subject:   "issuer-user-id",
		Email:     "person@example.com",
		Name:      "Person Example",
		Groups:    []string{"speedtest-admins"},
		ExpiresAt: expiresAt,
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("result = %#v, want %#v", result, want)
	}
}

func TestInvokeAuthenticateDoesNotStoreInvalidResponse(t *testing.T) {
	client := &fakeStateClient{}
	auth := &Auth{client: client}

	_, err := auth.invokeAuthenticate(context.Background(), envelope{value: OIDCResponse{}})
	if err == nil {
		t.Fatal("invokeAuthenticate() error = nil, want error")
	}
	if client.setCalled {
		t.Error("SetState was called for an invalid issuer response")
	}
}

func TestInvokeAuthenticateReturnsStoreError(t *testing.T) {
	storeErr := errors.New("store unavailable")
	client := &fakeStateClient{setErr: storeErr}
	auth := &Auth{client: client}

	_, err := auth.invokeAuthenticate(context.Background(), envelope{value: OIDCResponse{Subject: "issuer-user-id"}})
	if !errors.Is(err, storeErr) {
		t.Errorf("invokeAuthenticate() error = %v, want %v", err, storeErr)
	}
}

func TestInvokeBeginLoginStoresTransaction(t *testing.T) {
	client := &fakeStateClient{}
	authorizer := &authorizationURLBuilderStub{url: "https://issuer.example/authorize"}
	auth := &Auth{client: client, authorizer: authorizer}

	response, err := auth.invokeBeginLogin(context.Background())
	if err != nil {
		t.Fatalf("invokeBeginLogin() error = %v", err)
	}
	if response.RedirectURL != authorizer.url {
		t.Errorf("redirect URL = %q, want %q", response.RedirectURL, authorizer.url)
	}
	transaction := client.state.LoginTransaction
	if transaction == nil || transaction.State == "" || transaction.Verifier == "" || transaction.Nonce == "" {
		t.Errorf("stored transaction = %#v, want generated state, verifier, and nonce", transaction)
	}
	if transaction != nil && !transaction.ExpiresAt.After(time.Now()) {
		t.Errorf("transaction expiry = %v, want future time", transaction.ExpiresAt)
	}
	if !authorizer.called {
		t.Error("authorization URL builder was not called")
	}
	if !client.setCalled || client.setOpts == nil || client.setOpts.TTL != loginTransactionTTL {
		t.Errorf("transaction was not stored with the expected TTL: %#v", client.setOpts)
	}
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

	result, err := auth.invokeCompleteLogin(context.Background(), envelope{value: CompleteLoginRequest{
		Code: "authorization-code", State: "state",
	}})
	if err != nil {
		t.Fatalf("invokeCompleteLogin() error = %v", err)
	}
	if !exchanger.called || exchanger.code != "authorization-code" || exchanger.verifier != "verifier" || exchanger.nonce != "nonce" {
		t.Errorf("exchanger input = %#v", exchanger)
	}
	if client.state.LoginTransaction != nil {
		t.Error("login transaction was not cleared after successful authentication")
	}
	if !reflect.DeepEqual(client.state.OIDCResponse, response) {
		t.Errorf("stored response = %#v, want %#v", client.state.OIDCResponse, response)
	}
	if result.Subject != response.Subject || result.ExpiresAt != response.ExpiresAt {
		t.Errorf("result = %#v, want response summary", result)
	}
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
	if err == nil {
		t.Fatal("invokeCompleteLogin() error = nil, want state validation error")
	}
	if exchanger.called || client.setCalled {
		t.Error("state mismatch exchanged or mutated the login transaction")
	}
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
	if err != nil {
		t.Fatalf("invokeSession() error = %v", err)
	}

	want := AuthenticateResult{
		Subject:   "issuer-user-id",
		Email:     "person@example.com",
		Name:      "Person Example",
		Groups:    []string{"speedtest-admins"},
		ExpiresAt: expiresAt,
	}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("result = %#v, want %#v", result, want)
	}
	if client.haltCalled || client.deleteCalled {
		t.Error("valid session was halted or deleted")
	}
}

func TestInvokeSessionRejectsStateWithoutAuthenticatedSession(t *testing.T) {
	client := &fakeStateClient{}
	auth := &Auth{client: client}

	_, err := auth.invokeSession(context.Background())
	if !errors.Is(err, errNoAuthenticatedSession) {
		t.Errorf("invokeSession() error = %v, want %v", err, errNoAuthenticatedSession)
	}
	if client.haltCalled || client.deleteCalled {
		t.Error("empty state was halted or deleted")
	}
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
	if err != nil {
		t.Fatalf("invokeSession() error = %v", err)
	}
	if !refresher.called {
		t.Error("refresher was not called for an expiring session")
	}
	if client.state.OIDCResponse.AccessToken != "new-access-token" {
		t.Errorf("stored access token = %q, want refreshed token", client.state.OIDCResponse.AccessToken)
	}
	if !client.setCalled || client.setOpts == nil || client.setOpts.TTL <= 0 {
		t.Errorf("refreshed session was not persisted with a positive TTL: %#v", client.setOpts)
	}
	if got := result.(AuthenticateResult); got.ExpiresAt != newExpiry {
		t.Errorf("session expiry = %v, want %v", got.ExpiresAt, newExpiry)
	}
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
	if !errors.Is(err, refreshErr) {
		t.Errorf("invokeSession() error = %v, want %v", err, refreshErr)
	}
	if !client.haltCalled || !client.deleteCalled {
		t.Error("expired session was not halted and deleted after refresh failed")
	}
}

func TestInvokeSessionExpiresAndDeletesState(t *testing.T) {
	client := &fakeStateClient{state: authState{OIDCResponse: OIDCResponse{
		Subject:   "issuer-user-id",
		ExpiresAt: time.Now().Add(-time.Hour),
	}}}
	auth := &Auth{client: client}

	_, err := auth.invokeSession(context.Background())
	if err == nil || err.Error() != "expired session" {
		t.Errorf("invokeSession() error = %v, want expired session", err)
	}
	if !client.haltCalled {
		t.Error("Halt was not called for an expired session")
	}
	if !client.deleteCalled {
		t.Error("DeleteState was not called for an expired session")
	}
}

func TestInvokeSessionReturnsStateReadError(t *testing.T) {
	readErr := errors.New("state unavailable")
	auth := &Auth{client: &fakeStateClient{getErr: readErr}}

	_, err := auth.invokeSession(context.Background())
	if !errors.Is(err, readErr) {
		t.Errorf("invokeSession() error = %v, want %v", err, readErr)
	}
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

func (f *fakeStateClient) SetState(_ context.Context, state authState, opts *actor.SetStateOpts) error {
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

func (a *authorizationURLBuilderStub) AuthorizationURL(_ context.Context, _ LoginTransaction) (string, error) {
	a.called = true
	return a.url, a.err
}

func (e *tokenExchangerStub) Exchange(_ context.Context, code, verifier, nonce string) (OIDCResponse, error) {
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
