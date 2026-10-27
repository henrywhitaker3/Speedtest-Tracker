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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	assert.Equal(t, "new@example.com", response.Email)
	assert.Equal(t, "New Name", response.Name)
	assert.Equal(t, []string{"users"}, response.Groups)
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
	require.NoError(t, err)

	got, err := url.Parse(redirectURL)
	require.NoError(t, err)
	assert.Equal(t, issuer.server.URL+"/authorize", got.Scheme+"://"+got.Host+got.Path)
	query := got.Query()
	assert.Equal(t, "client-id", query.Get("client_id"))
	assert.Equal(t, builder.callbackURL, query.Get("redirect_uri"))
	assert.Equal(t, "openid profile email", query.Get("scope"))
	assert.Equal(t, "state", query.Get("state"))
	assert.Equal(t, "nonce", query.Get("nonce"))
	assert.Equal(t, "S256", query.Get("code_challenge_method"))
	assert.NotEmpty(t, query.Get("code_challenge"))
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
	require.NoError(t, err)
	assert.Equal(t, "access-token", response.AccessToken)
	assert.Equal(t, "refresh-token", response.RefreshToken)
	assert.Equal(t, "issuer-user-id", response.Subject)
	assert.Equal(t, "person@example.com", response.Email)
	assert.Equal(t, "Person Example", response.Name)
	assert.Equal(t, []string{"administrators", "users"}, response.Groups)
	assert.Equal(t, "openid profile email groups offline_access", response.Scope)
	assert.True(t, response.ExpiresAt.After(time.Now()))
	require.Len(t, issuer.tokenRequests, 1)
	request := issuer.tokenRequests[0]
	assert.Equal(t, "authorization_code", request.Get("grant_type"))
	assert.Equal(t, "code", request.Get("code"))
	assert.Equal(t, "verifier", request.Get("code_verifier"))
	assert.Equal(t, exchanger.callbackURL, request.Get("redirect_uri"))
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
		config: config.Auth{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			GroupsClaim:  "groups",
		},
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
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", refreshed.AccessToken)
	assert.Equal(t, "new-refresh-token", refreshed.RefreshToken)
	assert.Equal(t, "new@example.com", refreshed.Email)
	assert.Equal(t, "New Name", refreshed.Name)
	assert.Equal(t, []string{"users"}, refreshed.Groups)
	assert.Equal(t, "openid email", refreshed.Scope)
	assert.True(t, refreshed.ExpiresAt.After(time.Now()))
	require.Len(t, issuer.tokenRequests, 1)
	request := issuer.tokenRequests[0]
	assert.Equal(t, "refresh_token", request.Get("grant_type"))
	assert.Equal(t, "old-refresh-token", request.Get("refresh_token"))
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
	require.NoError(t, err)
	issuer := &oidcTestIssuer{key: key}
	issuer.server = httptest.NewServer(http.HandlerFunc(issuer.serveHTTP))
	t.Cleanup(issuer.server.Close)
	return issuer
}

func (i *oidcTestIssuer) provider(t *testing.T) *oidc.Provider {
	t.Helper()

	provider, err := oidc.NewProvider(context.Background(), i.server.URL)
	require.NoError(t, err)
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
		writeOIDCTestJSON(w, map[string]any{"keys": []map[string]string{
			{
				"kty": "RSA",
				"kid": "test-key",
				"use": "sig",
				"alg": "RS256",
				"n":   base64.RawURLEncoding.EncodeToString(i.key.PublicKey.N.Bytes()),
				"e": base64.RawURLEncoding.EncodeToString(
					big.NewInt(int64(i.key.PublicKey.E)).Bytes(),
				),
			},
		}})
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
	require.NoError(t, err)
	payload, err := json.Marshal(claims)
	require.NoError(t, err)
	signingInput := base64.RawURLEncoding.EncodeToString(
		header,
	) + "." + base64.RawURLEncoding.EncodeToString(
		payload,
	)
	sum := sha256.Sum256([]byte(signingInput))
	signature, err := rsa.SignPKCS1v15(rand.Reader, i.key, crypto.SHA256, sum[:])
	require.NoError(t, err)
	return signingInput + "." + base64.RawURLEncoding.EncodeToString(signature)
}

func writeOIDCTestJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}
