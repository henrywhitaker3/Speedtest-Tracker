package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"github.com/henrywhitaker3/windowframe/v2/http/handlers"
	"github.com/labstack/echo/v5"
	"golang.org/x/oauth2"
)

const (
	oidcStateCookieName    = "__Host-speedtest-oidc-state"
	oidcVerifierCookieName = "__Host-speedtest-oidc-verifier"
	oidcNonceCookieName    = "__Host-speedtest-oidc-nonce"
	oidcTransactionTTL     = 10 * time.Minute
)

type loginHandler struct {
	conf        config.Auth
	callbackURL string
	authorizer  authorizationURLBuilder
}

type authorizationURLBuilder interface {
	AuthorizationURL(context.Context, string, string, string) (string, error)
}

type oidcAuthorizer struct {
	conf        config.Auth
	callbackURL string
}

func (a oidcAuthorizer) AuthorizationURL(
	ctx context.Context,
	state, verifier, nonce string,
) (string, error) {
	provider, err := oidc.NewProvider(ctx, a.conf.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("discover OIDC provider: %w", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     a.conf.ClientID,
		ClientSecret: string(a.conf.ClientSecret),
		RedirectURL:  a.callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       strings.Fields(a.conf.Claims),
	}

	return oauthConfig.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("nonce", nonce),
	), nil
}

func (l *loginHandler) Handler() handlers.Handler[any, any] {
	return func(c *echo.Context, _ any) (*any, error) {
		sessionID, err := randomValue(32)
		if err != nil {
			return nil, fmt.Errorf("generate session ID: %w", err)
		}
		state, err := randomValue(32)
		if err != nil {
			return nil, fmt.Errorf("generate OIDC state: %w", err)
		}
		verifier, err := randomValue(32)
		if err != nil {
			return nil, fmt.Errorf("generate PKCE verifier: %w", err)
		}
		nonce, err := randomValue(32)
		if err != nil {
			return nil, fmt.Errorf("generate OIDC nonce: %w", err)
		}

		authorizer := l.authorizer
		if authorizer == nil {
			authorizer = oidcAuthorizer{
				conf:        l.conf,
				callbackURL: l.callbackURL,
			}
		}
		authorizationURL, err := authorizer.AuthorizationURL(
			c.Request().Context(), state, verifier, nonce,
		)
		if err != nil {
			return nil, fmt.Errorf("build OIDC authorization URL: %w", err)
		}

		setSessionCookie(c, sessionCookieName, sessionID, 0)
		setSessionCookie(c, oidcStateCookieName, state, oidcTransactionTTL)
		setSessionCookie(c, oidcVerifierCookieName, verifier, oidcTransactionTTL)
		setSessionCookie(c, oidcNonceCookieName, nonce, oidcTransactionTTL)

		if err := c.Redirect(http.StatusFound, authorizationURL); err != nil {
			return nil, fmt.Errorf("redirect to OIDC issuer: %w", err)
		}
		return nil, nil
	}
}

func randomValue(bytes int) (string, error) {
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func setSessionCookie(c *echo.Context, name, value string, ttl time.Duration) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if ttl > 0 {
		cookie.MaxAge = int(ttl.Seconds())
		cookie.Expires = time.Now().Add(ttl)
	}
	c.SetCookie(cookie)
}

func (l *loginHandler) Middleware() []echo.MiddlewareFunc {
	return nil
}

func (l *loginHandler) Metadata() handlers.Metadata {
	return handlers.Metadata{
		Name:         "Auth login",
		Tag:          "Auth",
		Path:         "/v1/auth/login",
		Code:         http.StatusOK,
		Method:       http.MethodGet,
		GenerateSpec: false,
	}
}
