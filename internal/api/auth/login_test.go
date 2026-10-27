package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authactor "github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/auth"
	"github.com/labstack/echo/v5"
)

func TestLoginRedirectsToActorProvidedURL(t *testing.T) {
	service := &actorInvokerStub{envelope: loginEnvelope{value: authactor.LoginResponse{
		RedirectURL: "https://issuer.example/authorize",
	}}}
	handler := (&loginHandler{service: service}).Handler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-id"})
	ctx := echo.NewContext(req, recorder)

	response, err := handler(ctx, nil)
	require.Equalf(t, nil, err, "Handler() error = %v", err)
	assert.Nilf(t, response, "Handler() response = %#v, want nil", response)
	assert.Equalf(
		t,
		http.StatusFound,
		recorder.Code,
		"status = %d, want %d",
		recorder.Code,
		http.StatusFound,
	)
	assert.Equal(t, "https://issuer.example/authorize", recorder.Header().Get(echo.HeaderLocation))
	assert.Equal(t, authactor.Actor, service.actorType)
	assert.Equal(t, authactor.BeginLogin, service.method)
	assert.Equal(t, "session-id", service.actorID)
}

func TestLoginWrapsActorError(t *testing.T) {
	actorErr := errors.New("actor unavailable")
	handler := (&loginHandler{service: &actorInvokerStub{err: actorErr}}).Handler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-id"})
	ctx := echo.NewContext(req, recorder)

	_, err := handler(ctx, nil)
	assert.ErrorIs(t, err, actorErr)
}

type loginEnvelope struct {
	value authactor.LoginResponse
}

func (e loginEnvelope) Decode(into any) error {
	*into.(*authactor.LoginResponse) = e.value
	return nil
}
