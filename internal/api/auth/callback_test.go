package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/italypaleale/francis/actor"
	"github.com/labstack/echo/v5"
)

func TestCallbackInvokesAuthenticateForCookieSession(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&callbackHandler{
		service:            service,
		successRedirectURL: "https://speedtest.example/",
	}).Handler()

	ctx := callbackContext(t, "session-id")
	got, err := handler(
		ctx,
		CallbackRequest{Code: "authorization-code", State: "state"},
	)
	if err != nil {
		t.Fatalf("Handler() error = %v", err)
	}
	if got != nil {
		t.Errorf("Handler() result = %#v, want nil", got)
	}
	if location := ctx.Response().Header().Get(echo.HeaderLocation); location != "https://speedtest.example/" {
		t.Errorf("Location = %q, want dashboard", location)
	}
	if service.actorType != authactor.Actor {
		t.Errorf("actor type = %q, want %q", service.actorType, authactor.Actor)
	}
	if service.actorID != "session-id" {
		t.Errorf("actor ID = %q, want session cookie value", service.actorID)
	}
	if service.method != authactor.CompleteLogin {
		t.Errorf("actor method = %q, want %q", service.method, authactor.CompleteLogin)
	}
	want := authactor.CompleteLoginRequest{Code: "authorization-code", State: "state"}
	if !reflect.DeepEqual(service.data, want) {
		t.Errorf("actor data = %#v, want %#v", service.data, want)
	}
}

func TestCallbackRejectsMissingSessionCookie(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, ""),
		CallbackRequest{Code: "authorization-code", State: "state"},
	)
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusUnauthorized {
		t.Errorf("Handler() error = %v, want 401", err)
	}
	if service.invoked {
		t.Error("actor was invoked without a session cookie")
	}
}

func TestCallbackWrapsActorError(t *testing.T) {
	invokeErr := errors.New("actor unavailable")
	service := &actorInvokerStub{err: invokeErr}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, "session-id"),
		CallbackRequest{Code: "authorization-code", State: "state"},
	)
	if !errors.Is(err, invokeErr) {
		t.Errorf("Handler() error = %v, want %v", err, invokeErr)
	}
}

func TestCallbackReturnsIssuerErrorWithoutInvokingActor(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, "session-id"),
		CallbackRequest{Error: "access_denied"},
	)
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != http.StatusBadRequest {
		t.Errorf("Handler() error = %v, want 400", err)
	}
	if service.invoked {
		t.Error("callback invoked the actor after the issuer returned an error")
	}
}

func callbackContext(t *testing.T, sessionID string) *echo.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/callback", nil)
	if sessionID != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: sessionID})
	}
	return echo.NewContext(req, httptest.NewRecorder())
}

type actorInvokerStub struct {
	envelope  actor.Envelope
	err       error
	actorType string
	actorID   string
	method    string
	data      any
	invoked   bool
}

func (s *actorInvokerStub) Invoke(
	_ context.Context,
	actorType, actorID, method string,
	data any,
	_ ...actor.InvokeOption,
) (actor.Envelope, error) {
	s.invoked = true
	s.actorType = actorType
	s.actorID = actorID
	s.method = method
	s.data = data
	return s.envelope, s.err
}
