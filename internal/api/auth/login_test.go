package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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
	if err != nil {
		t.Fatalf("Handler() error = %v", err)
	}
	if response != nil {
		t.Errorf("Handler() response = %#v, want nil", response)
	}
	if recorder.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusFound)
	}
	if recorder.Header().Get(echo.HeaderLocation) != "https://issuer.example/authorize" {
		t.Errorf("Location = %q, want issuer URL", recorder.Header().Get(echo.HeaderLocation))
	}
	if service.actorType != authactor.Actor || service.method != authactor.BeginLogin || service.actorID != "session-id" {
		t.Errorf("actor invocation = type=%q ID=%q method=%q", service.actorType, service.actorID, service.method)
	}
}

func TestLoginWrapsActorError(t *testing.T) {
	actorErr := errors.New("actor unavailable")
	handler := (&loginHandler{service: &actorInvokerStub{err: actorErr}}).Handler()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil)
	req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: "session-id"})
	ctx := echo.NewContext(req, recorder)

	_, err := handler(ctx, nil)
	if !errors.Is(err, actorErr) {
		t.Errorf("Handler() error = %v, want %v", err, actorErr)
	}
}

type loginEnvelope struct {
	value authactor.LoginResponse
}

func (e loginEnvelope) Decode(into any) error {
	*into.(*authactor.LoginResponse) = e.value
	return nil
}
