package auth

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
)

func TestApplyIdentityClaimsReplacesRefreshedGroups(t *testing.T) {
	response := OIDCResponse{
		Email:  "old@example.com",
		Name:   "Old Name",
		Groups: []string{"speedtest-admins", "users"},
	}

	applyIdentityClaims(&response, oidcClaims{
		Email:  "new@example.com",
		Name:   "New Name",
		Groups: []string{"users"},
	})

	if response.Email != "new@example.com" || response.Name != "New Name" {
		t.Errorf("identity = email %q name %q, want refreshed claims", response.Email, response.Name)
	}
	if !reflect.DeepEqual(response.Groups, []string{"users"}) {
		t.Errorf("groups = %#v, want refreshed groups without admin access", response.Groups)
	}
}

func TestOIDCAuthorizationURLBuilderUsesConfiguredProvider(t *testing.T) {
	issuer := newOIDCTestIssuer(t)
	provider := issuer.provider(t)
	builder := oidcAuthorizationURLBuilder{
		config: config.Auth{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Scopes:       "openid profile email",
		},
		callbackURL: "https://speedtest.example/v1/auth/callback",
		provider:    provider,
	}

	redirectURL, err := builder.AuthorizationURL(context.Background(), LoginTransaction{
		State:    "state",
		Verifier: "verifier",
		Nonce:    "nonce",
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := url.Parse(redirectURL)
	if err != nil {
		t.Fatal(err)
	}
	if got.Scheme+"://"+got.Host+got.Path != issuer.server.URL+"/authorize" {
		t.Errorf("authorization endpoint = %q, want %q", got.Scheme+"://"+got.Host+got.Path, issuer.server.URL+"/authorize")
	}
	query := got.Query()
	if query.Get("client_id") != "client-id" || query.Get("redirect_uri") != builder.callbackURL {
		t.Errorf("authorization query = %q, want configured client and redirect URL", query)
	}
	if query.Get("scope") != "openid profile email" || query.Get("state") != "state" || query.Get("nonce") != "nonce" {
		t.Errorf("authorization query = %q, want configured scopes, state, and nonce", query)
	}
	if query.Get("code_challenge_method") != "S256" || query.Get("code_challenge") == "" {
		t.Errorf("authorization query = %q, want PKCE S256 challenge", query)
	}
}

func TestOIDCTokenExchangerExchangesAndVerifiesIDToken(t *testing.T) {
	issuer := newOIDCTestIssuer(t)
	issuer.tokenResponse = func(values url.Values) oidcTestTokenResponse {
		return oidcTestTokenResponse{
			accessToken:  "access-token",
			refreshToken: "refresh-token",
			idToken: issuer.signedIDToken(t, map[string]any{
				"sub":    "issuer-user-id",
				"nonce":  "nonce",
				"email":  "person@example.com",
				"name":   "Person Example",
				"groups": []string{"administrators", "users"},
			}),
			scope: "openid profile email groups offline_access",
		}
	}

	exchanger := oidcTokenExchanger{
		config: config.Auth{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			GroupsClaim:  "groups",
		},
		callbackURL: "https://speedtest.example/v1/auth/callback",
		provider:    issuer.provider(t),
	}

	response, err := exchanger.Exchange(context.Background(), "code", "verifier", "nonce")
	if err != nil {
		t.Fatal(err)
	}
	if response.AccessToken != "access-token" || response.RefreshToken != "refresh-token" || response.Subject != "issuer-user-id" {
		t.Errorf("response credentials = %#v", response)
	}
	if response.Email != "person@example.com" || response.Name != "Person Example" ||
		!reflect.DeepEqual(response.Groups, []string{"administrators", "users"}) {
		t.Errorf("response identity = %#v", response)
	}
	if response.Scope != "openid profile email groups offline_access" || !response.ExpiresAt.After(time.Now()) {
		t.Errorf("response scope or expiry = %#v", response)
	}
	if len(issuer.tokenRequests) != 1 {
		t.Fatalf("token requests = %d, want 1", len(issuer.tokenRequests))
	}
	request := issuer.tokenRequests[0]
	if request.Get("grant_type") != "authorization_code" || request.Get("code") != "code" ||
		request.Get("code_verifier") != "verifier" || request.Get("redirect_uri") != exchanger.callbackURL {
		t.Errorf("token request = %q, want authorization-code exchange request", request)
	}
}

func TestOIDCTokenRefresherRotatesTokensAndIdentity(t *testing.T) {
	issuer := newOIDCTestIssuer(t)
	issuer.tokenResponse = func(values url.Values) oidcTestTokenResponse {
		return oidcTestTokenResponse{
			accessToken:  "new-access-token",
			refreshToken: "new-refresh-token",
			idToken: issuer.signedIDToken(t, map[string]any{
				"sub":    "issuer-user-id",
				"email":  "new@example.com",
				"name":   "New Name",
				"groups": []string{"users"},
			}),
			scope: "openid email",
		}
	}

	refresher := oidcTokenRefresher{
		config:   config.Auth{ClientID: "client-id", ClientSecret: "client-secret", GroupsClaim: "groups"},
		provider: issuer.provider(t),
	}
	refreshed, err := refresher.Refresh(context.Background(), OIDCResponse{
		AccessToken:  "old-access-token",
		RefreshToken: "old-refresh-token",
		Subject:      "issuer-user-id",
		Email:        "old@example.com",
		Name:         "Old Name",
		Groups:       []string{"administrators", "users"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if refreshed.AccessToken != "new-access-token" || refreshed.RefreshToken != "new-refresh-token" {
		t.Errorf("refreshed credentials = %#v", refreshed)
	}
	if refreshed.Email != "new@example.com" || refreshed.Name != "New Name" ||
		!reflect.DeepEqual(refreshed.Groups, []string{"users"}) {
		t.Errorf("refreshed identity = %#v", refreshed)
	}
	if refreshed.Scope != "openid email" || !refreshed.ExpiresAt.After(time.Now()) {
		t.Errorf("refreshed scope or expiry = %#v", refreshed)
	}
	if len(issuer.tokenRequests) != 1 || issuer.tokenRequests[0].Get("grant_type") != "refresh_token" ||
		issuer.tokenRequests[0].Get("refresh_token") != "old-refresh-token" {
		t.Errorf("token requests = %#v, want refresh-token request", issuer.tokenRequests)
	}
}

type oidcTestIssuer struct {
	key           *rsa.PrivateKey
	server        *httptest.Server
	tokenRequests []url.Values
	tokenResponse func(url.Values) oidcTestTokenResponse
}

type oidcTestTokenResponse struct {
	accessToken  string
	refreshToken string
	idToken      string
	scope        string
}

func newOIDCTestIssuer(t *testing.T) *oidcTestIssuer {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	issuer := &oidcTestIssuer{key: key}
	issuer.server = httptest.NewServer(http.HandlerFunc(issuer.serveHTTP))
	t.Cleanup(issuer.server.Close)
	return issuer
}

func (i *oidcTestIssuer) provider(t *testing.T) *oidc.Provider {
	t.Helper()

	provider, err := oidc.NewProvider(context.Background(), i.server.URL)
	if err != nil {
		t.Fatal(err)
	}
	return provider
}

func (i *oidcTestIssuer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		writeOIDCTestJSON(w, map[string]string{
			"issuer":                 i.server.URL,
			"authorization_endpoint": i.server.URL + "/authorize",
			"token_endpoint":         i.server.URL + "/token",
			"jwks_uri":               i.server.URL + "/keys",
		})
	case "/keys":
		writeOIDCTestJSON(w, map[string]any{"keys": []map[string]string{{
			"kty": "RSA",
			"kid": "test-key",
			"use": "sig",
			"alg": "RS256",
			"n":   base64.RawURLEncoding.EncodeToString(i.key.PublicKey.N.Bytes()),
			"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(i.key.PublicKey.E)).Bytes()),
		}}})
	case "/token":
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		i.tokenRequests = append(i.tokenRequests, r.Form)
		if i.tokenResponse == nil {
			http.Error(w, "unexpected token request", http.StatusInternalServerError)
			return
		}
		response := i.tokenResponse(r.Form)
		writeOIDCTestJSON(w, map[string]any{
			"access_token":  response.accessToken,
			"refresh_token": response.refreshToken,
			"id_token":      response.idToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"scope":         response.scope,
		})
	default:
		http.NotFound(w, r)
	}
}

func (i *oidcTestIssuer) signedIDToken(t *testing.T, claims map[string]any) string {
	t.Helper()

	for key, value := range map[string]any{
		"iss": i.server.URL,
		"aud": "client-id",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	} {
		if _, ok := claims[key]; !ok {
			claims[key] = value
		}
	}
	header, err := json.Marshal(map[string]string{"alg": "RS256", "kid": "test-key", "typ": "JWT"})
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	signingInput := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	sum := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, i.key, crypto.SHA256, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func writeOIDCTestJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
