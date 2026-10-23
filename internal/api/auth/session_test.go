package auth

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

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
	if err != nil {
		t.Fatalf("Handler() error = %v", err)
	}
	if !reflect.DeepEqual(got, &want) {
		t.Errorf("Handler() result = %#v, want %#v", got, &want)
	}
	if service.actorType != authactor.Actor || service.actorID != "session-id" || service.method != authactor.Session {
		t.Errorf("actor invocation = type=%q ID=%q method=%q", service.actorType, service.actorID, service.method)
	}
}

func TestSessionRejectsMissingCookie(t *testing.T) {
	service := &actorInvokerStub{}
	handler := (&sessionHandler{service: service}).Handler()

	_, err := handler(callbackContext(t, ""), nil)
	assertHTTPError(t, err, http.StatusUnauthorized)
	if service.invoked {
		t.Error("actor was invoked without a session cookie")
	}
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
		t.Errorf("error = %v, want HTTP %d", err, wantCode)
	}
}
