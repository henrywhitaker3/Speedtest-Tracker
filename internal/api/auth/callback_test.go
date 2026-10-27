package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.Equalf(t, nil, err, "Handler() error = %v", err)
	assert.Nilf(t, got, "Handler() result = %#v, want nil", got)
	assert.Equal(t, "https://speedtest.example/", ctx.Response().Header().Get(echo.HeaderLocation))
	assert.Equalf(
		t,
		authactor.Actor,
		service.actorType,
		"actor type = %q, want %q",
		service.actorType,
		authactor.Actor,
	)
	assert.Equal(t, "session-id", service.actorID)
	assert.Equalf(
		t,
		authactor.CompleteLogin,
		service.method,
		"actor method = %q, want %q",
		service.method,
		authactor.CompleteLogin,
	)
	require.Equal(
		t,
		authactor.CompleteLoginRequest{Code: "authorization-code", State: "state"},
		service.data,
	)
}

func TestCallbackRejectsMissingSessionCookie(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, ""),
		CallbackRequest{Code: "authorization-code", State: "state"},
	)
	var httpErr *echo.HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
	assert.False(t, service.invoked)
}

func TestCallbackWrapsActorError(t *testing.T) {
	invokeErr := errors.New("actor unavailable")
	service := &actorInvokerStub{err: invokeErr}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, "session-id"),
		CallbackRequest{Code: "authorization-code", State: "state"},
	)
	assert.ErrorIs(t, err, invokeErr)
}

func TestCallbackReturnsIssuerErrorWithoutInvokingActor(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&callbackHandler{service: service}).Handler()

	_, err := handler(
		callbackContext(t, "session-id"),
		CallbackRequest{Error: "access_denied"},
	)
	var httpErr *echo.HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	assert.False(t, service.invoked)
}

func callbackContext(t *testing.T, sessionID string) *echo.Context {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/callback", nil)
	if sessionID != "" {
		req.AddCookie(&http.Cookie{Name: SessionCookieName, Value: sessionID})
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
