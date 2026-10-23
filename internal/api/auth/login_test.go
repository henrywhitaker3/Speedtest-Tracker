package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/labstack/echo/v5"
)

func TestLoginRedirectsToIssuerAndCreatesTransactionCookies(t *testing.T) {
	authorizer := &authorizerStub{url: "https://issuer.example/authorize"}
	handler := (&loginHandler{
		conf:        config.Auth{ClientID: "speedtest-tracker"},
		callbackURL: "https://speedtest.example/v1/auth/callback",
		authorizer:  authorizer,
	}).Handler()
	recorder := httptest.NewRecorder()
	ctx := echo.NewContext(httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil), recorder)

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
	if recorder.Header().Get(echo.HeaderLocation) != authorizer.url {
		t.Errorf("Location = %q, want %q", recorder.Header().Get(echo.HeaderLocation), authorizer.url)
	}
	if authorizer.state == "" || authorizer.verifier == "" || authorizer.nonce == "" {
		t.Errorf("OIDC request values = state=%q verifier=%q nonce=%q, want non-empty values", authorizer.state, authorizer.verifier, authorizer.nonce)
	}

	cookies := cookiesByName(recorder.Result().Cookies())
	for _, name := range []string{sessionCookieName, oidcStateCookieName, oidcVerifierCookieName, oidcNonceCookieName} {
		cookie, ok := cookies[name]
		if !ok {
			t.Errorf("missing %s cookie", name)
			continue
		}
		if !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" {
			t.Errorf("cookie %s has insecure attributes: %#v", name, cookie)
		}
	}
	if cookies[sessionCookieName].Value == "" {
		t.Error("session cookie is empty")
	}
	if cookies[oidcStateCookieName].Value != authorizer.state {
		t.Error("state cookie does not match authorization request")
	}
	if cookies[oidcVerifierCookieName].Value != authorizer.verifier {
		t.Error("verifier cookie does not match authorization request")
	}
	if cookies[oidcNonceCookieName].Value != authorizer.nonce {
		t.Error("nonce cookie does not match authorization request")
	}
	if cookies[oidcStateCookieName].MaxAge != int(oidcTransactionTTL.Seconds()) {
		t.Errorf("state cookie MaxAge = %d, want %d", cookies[oidcStateCookieName].MaxAge, int(oidcTransactionTTL.Seconds()))
	}
}

func TestLoginDoesNotSetCookiesWhenIssuerDiscoveryFails(t *testing.T) {
	issuerErr := errors.New("issuer unavailable")
	handler := (&loginHandler{authorizer: &authorizerStub{err: issuerErr}}).Handler()
	recorder := httptest.NewRecorder()
	ctx := echo.NewContext(httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil), recorder)

	_, err := handler(ctx, nil)
	if !errors.Is(err, issuerErr) {
		t.Errorf("Handler() error = %v, want %v", err, issuerErr)
	}
	if len(recorder.Result().Cookies()) != 0 {
		t.Errorf("cookies = %#v, want none", recorder.Result().Cookies())
	}
}

func TestOIDCAuthorizerBuildsAuthorizationCodePKCEURL(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/openid-configuration" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		_, _ = w.Write([]byte(`{"issuer":"` + server.URL + `","authorization_endpoint":"https://issuer.example/authorize","token_endpoint":"https://issuer.example/token","jwks_uri":"https://issuer.example/keys"}`))
	}))
	defer server.Close()

	authorizer := oidcAuthorizer{
		conf: config.Auth{
			IssuerURL: server.URL,
			ClientID:  "speedtest-tracker",
			Claims:    "openid profile email",
		},
		callbackURL: "https://speedtest.example/v1/auth/callback",
	}
	got, err := authorizer.AuthorizationURL(context.Background(), "state", "verifier-value", "nonce")
	if err != nil {
		t.Fatalf("AuthorizationURL() error = %v", err)
	}
	url, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}
	query := url.Query()
	if query.Get("response_type") != "code" || query.Get("client_id") != "speedtest-tracker" || query.Get("redirect_uri") != authorizer.callbackURL {
		t.Errorf("authorization request = %q", got)
	}
	if query.Get("state") != "state" || query.Get("nonce") != "nonce" {
		t.Errorf("authorization request is missing state or nonce: %q", got)
	}
	if query.Get("code_challenge") == "" || query.Get("code_challenge_method") != "S256" {
		t.Errorf("authorization request is missing S256 PKCE: %q", got)
	}
}

func cookiesByName(cookies []*http.Cookie) map[string]*http.Cookie {
	result := make(map[string]*http.Cookie, len(cookies))
	for _, cookie := range cookies {
		result[cookie.Name] = cookie
	}
	return result
}

type authorizerStub struct {
	url      string
	err      error
	state    string
	verifier string
	nonce    string
}

func (a *authorizerStub) AuthorizationURL(_ context.Context, state, verifier, nonce string) (string, error) {
	a.state = state
	a.verifier = verifier
	a.nonce = nonce
	return a.url, a.err
}
