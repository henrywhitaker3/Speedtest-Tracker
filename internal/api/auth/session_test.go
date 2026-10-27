package auth

import (
	"errors"
	"net/http"
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/labstack/echo/v5"
)

func TestSessionReturnsAuthenticatedUser(t *testing.T) {
	expiresAt := time.Now().Add(time.Hour).UTC()
	want := authactor.AuthenticateResult{
		Subject:   "issuer-user-id",
		Email:     "person@example.com",
		Groups:    []string{"speedtest-admins"},
		ExpiresAt: expiresAt,
	}
	service := &actorInvokerStub{envelope: sessionEnvelope{value: want}}
	handler := (&sessionHandler{service: service}).Handler()

	got, err := handler(callbackContext(t, "session-id"), nil)
	require.Equalf(t, nil, err, "Handler() error = %v", err)
	assert.Equal(t, &want, got)
	assert.Equal(t, authactor.Actor, service.actorType)
	assert.Equal(t, authactor.Session, service.method)
	assert.Equal(t, "session-id", service.actorID)
}

func TestSessionRejectsMissingCookie(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&sessionHandler{service: service}).Handler()

	_, err := handler(callbackContext(t, ""), nil)
	assertHTTPError(t, err, http.StatusUnauthorized)
	assert.False(t, service.invoked)
}

type sessionEnvelope struct {
	value authactor.AuthenticateResult
}

func (e sessionEnvelope) Decode(into any) error {
	*into.(*authactor.AuthenticateResult) = e.value
	return nil
}

func assertHTTPError(t *testing.T, err error, wantCode int) {
	t.Helper()
	var httpErr *echo.HTTPError
	if !errors.As(err, &httpErr) || httpErr.Code != wantCode {
		assert.Failf(t, "assertion failed", "error = %v, want HTTP %d", err, wantCode)
	}
}
