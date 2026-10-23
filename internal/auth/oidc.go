package auth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/henrywhitaker3/Speedtest-Tracker/v2/internal/config"
	"golang.org/x/oauth2"
)

type tokenRefresher interface {
	Refresh(context.Context, OIDCResponse) (OIDCResponse, error)
}

type tokenExchanger interface {
	Exchange(context.Context, string, string, string) (OIDCResponse, error)
}

type authorizationURLBuilder interface {
	AuthorizationURL(context.Context, LoginTransaction) (string, error)
}

type oidcAuthorizationURLBuilder struct {
	config      config.Auth
	callbackURL string
}

func (a oidcAuthorizationURLBuilder) AuthorizationURL(
	ctx context.Context,
	transaction LoginTransaction,
) (string, error) {
	provider, err := oidc.NewProvider(ctx, a.config.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("discover OIDC provider: %w", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     a.config.ClientID,
		ClientSecret: string(a.config.ClientSecret),
		RedirectURL:  a.callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       strings.Fields(a.config.Scopes),
	}
	return oauthConfig.AuthCodeURL(
		transaction.State,
		oauth2.S256ChallengeOption(transaction.Verifier),
		oauth2.SetAuthURLParam("nonce", transaction.Nonce),
	), nil
}

type oidcTokenExchanger struct {
	config      config.Auth
	callbackURL string
}

func (e oidcTokenExchanger) Exchange(
	ctx context.Context,
	code, verifier, nonce string,
) (OIDCResponse, error) {
	provider, err := oidc.NewProvider(ctx, e.config.IssuerURL)
	if err != nil {
		return OIDCResponse{}, fmt.Errorf("discover OIDC provider: %w", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     e.config.ClientID,
		ClientSecret: string(e.config.ClientSecret),
		RedirectURL:  e.callbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       strings.Fields(e.config.Scopes),
	}
	token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return OIDCResponse{}, fmt.Errorf("exchange authorization code: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return OIDCResponse{}, errors.New("OIDC token response has no ID token")
	}
	idToken, err := provider.VerifierContext(ctx, &oidc.Config{ClientID: e.config.ClientID}).
		Verify(ctx, rawIDToken)
	if err != nil {
		return OIDCResponse{}, fmt.Errorf("verify ID token: %w", err)
	}
	if subtle.ConstantTimeCompare([]byte(idToken.Nonce), []byte(nonce)) != 1 {
		return OIDCResponse{}, errors.New("verify ID token: nonce mismatch")
	}

	claims, err := extractClaims(idToken, e.config.GroupsClaim)
	if err != nil {
		return OIDCResponse{}, err
	}
	scope, _ := token.Extra("scope").(string)
	return OIDCResponse{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		IDToken:      rawIDToken,
		Scope:        scope,
		ExpiresAt:    idToken.Expiry,
		Subject:      idToken.Subject,
		Email:        claims.Email,
		Name:         claims.Name,
		Groups:       claims.Groups,
	}, nil
}

type oidcClaims struct {
	Email   string
	Name    string
	Picture string
	Groups  []string
}

func extractClaims(idToken *oidc.IDToken, groupsClaim string) (oidcClaims, error) {
	var claims struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return oidcClaims{}, fmt.Errorf("decode ID token claims: %w", err)
	}

	var rawClaims map[string]json.RawMessage
	if err := idToken.Claims(&rawClaims); err != nil {
		return oidcClaims{}, fmt.Errorf("decode ID token group claim: %w", err)
	}

	result := oidcClaims{Email: claims.Email, Name: claims.Name}
	groups, ok := rawClaims[groupsClaim]
	if !ok {
		return result, nil
	}
	if err := json.Unmarshal(groups, &result.Groups); err != nil {
		return oidcClaims{}, fmt.Errorf("decode ID token group claim %q: %w", groupsClaim, err)
	}
	return result, nil
}

type oidcTokenRefresher struct {
	config config.Auth
}

func (r oidcTokenRefresher) Refresh(
	ctx context.Context,
	current OIDCResponse,
) (OIDCResponse, error) {
	if current.RefreshToken == "" {
		return OIDCResponse{}, errors.New("no refresh token")
	}

	provider, err := oidc.NewProvider(ctx, r.config.IssuerURL)
	if err != nil {
		return OIDCResponse{}, fmt.Errorf("discover OIDC provider: %w", err)
	}
	oauthConfig := oauth2.Config{
		ClientID:     r.config.ClientID,
		ClientSecret: string(r.config.ClientSecret),
		Endpoint:     provider.Endpoint(),
	}
	token, err := oauthConfig.TokenSource(ctx, &oauth2.Token{
		RefreshToken: current.RefreshToken,
	}).Token()
	if err != nil {
		return OIDCResponse{}, fmt.Errorf("refresh token: %w", err)
	}

	refreshed := current
	refreshed.AccessToken = token.AccessToken
	refreshed.TokenType = token.TokenType
	if token.RefreshToken != "" {
		refreshed.RefreshToken = token.RefreshToken
	}
	if scope, ok := token.Extra("scope").(string); ok && scope != "" {
		refreshed.Scope = scope
	}

	if rawIDToken, ok := token.Extra("id_token").(string); ok && rawIDToken != "" {
		idToken, err := provider.VerifierContext(ctx, &oidc.Config{ClientID: r.config.ClientID}).
			Verify(ctx, rawIDToken)
		if err != nil {
			return OIDCResponse{}, fmt.Errorf("verify refreshed ID token: %w", err)
		}
		if idToken.Subject != current.Subject {
			return OIDCResponse{}, errors.New("refreshed ID token subject does not match session")
		}
		claims, err := extractClaims(idToken, r.config.GroupsClaim)
		if err != nil {
			return OIDCResponse{}, fmt.Errorf("decode refreshed ID token claims: %w", err)
		}
		refreshed.IDToken = rawIDToken
		refreshed.ExpiresAt = idToken.Expiry
		applyIdentityClaims(&refreshed, claims)
	} else {
		refreshed.ExpiresAt = token.Expiry
	}
	if refreshed.ExpiresAt.IsZero() {
		return OIDCResponse{}, errors.New("refreshed token has no expiry")
	}

	return refreshed, nil
}

func applyIdentityClaims(response *OIDCResponse, claims oidcClaims) {
	response.Email = claims.Email
	response.Name = claims.Name
	response.Groups = claims.Groups
}
