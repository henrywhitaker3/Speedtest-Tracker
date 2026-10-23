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
	haltCalled   bool
	deleteCalled bool
}

func (f *fakeStateClient) GetState(context.Context) (authState, error) {
	return f.state, f.getErr
}

func (f *fakeStateClient) SetState(_ context.Context, state authState, _ *actor.SetStateOpts) error {
	f.setCalled = true
	if f.setErr != nil {
		return f.setErr
	}
	f.state = state
	return nil
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
	response, ok := into.(*OIDCResponse)
	if !ok {
		return errors.New("unexpected decode target")
	}
	*response = e.value.(OIDCResponse)
	return nil
}
