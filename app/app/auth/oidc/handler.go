package oidc

import (
	"context"
	"net/http"
	"raumzeitalpaka/app/auth"
	"raumzeitalpaka/config"
	"raumzeitalpaka/ports/www/render"
	"strconv"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

var (
	CallbackPath = "/login/callback"
)

func NewHandler(ctx context.Context, cfg config.Config, syncer Syncer) (r chi.Router, err error) {
	log := zap.L()
	redirectURI := config.BaseUrl() + CallbackPath

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}
	var verifier = provider.Verifier(&oidc.Config{ClientID: cfg.ClientID})

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  redirectURI,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	// State generator for CSRF protection
	state := func() string {
		return uuid.New().String()
	}

	var handleRedirect http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, oauth2Config.AuthCodeURL(state()), http.StatusFound)
	})

	var handleOAuth2Callback http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify state and errors.

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			render.Error(log, w, http.StatusUnauthorized, "oauth2 exchange failed", err)
			return
		}

		// Extract the ID Token from OAuth2 token.
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			render.Error(log, w, http.StatusUnauthorized, "missing oauth2 token", err)
			return
		}

		// Parse and verify ID Token payload.
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			render.Error(log, w, http.StatusUnauthorized, "failed to verify ID token", err)
			return
		}

		// Extract custom claims
		var claims claims
		//var claims = make(map[string]any)
		if err := idToken.Claims(&claims); err != nil {
			render.Error(log, w, http.StatusUnauthorized, "failed to parse ID Token claims", err)
			return
		}
		log.Info("claims", zap.Any("claims", claims))

		userId, err := strconv.ParseInt(claims.UserId, 10, 64)
		if err != nil {
			render.Error(log, w, http.StatusBadRequest, "failed to parse userId", err)
			return
		}

		err = syncer.Sync(int(userId), claims.Username, claims.Groups)
		if err != nil {
			render.Error(log, w, http.StatusInternalServerError, "failed to sync user", err)
			return
		}

		jwt, err := auth.NewJWT(auth.Claims{
			UserId: int(userId),
		})
		auth.SetSessionCookie(w, jwt)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	mux := chi.NewMux()
	mux.Handle("/login", handleRedirect)
	mux.Handle(CallbackPath, handleOAuth2Callback)
	mux.Handle("/logout", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		auth.ClearSessionCookie(writer)
		writer.Write([]byte("logged out"))
	}))
	return mux, nil
}

type claims struct {
	Email    string   `json:"email"`
	Username string   `json:"preferred_username"`
	Groups   []string `json:"groups"`
	UserId   string   `json:"sub"`
}
