package oauth2

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"golang.org/x/oauth2"
)

type Config struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
}

type Service struct {
	config   *oauth2.Config
	verifier *oidc.IDTokenVerifier

	stateStorage        *ttlcache.Cache[string, string]
	codeVerifierStorage *ttlcache.Cache[string, string]
}

func NewService(config Config, redirectURL string) (*Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: redirectURL,
		Scopes:      []string{"openid", "profile"},
	}

	provider, err := oidc.NewProvider(context.Background(), config.Issuer)
	if err != nil {
		return nil, fmt.Errorf("new oidc provider: %w", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	stateStorage := ttlcache.New(
		ttlcache.WithTTL[string, string](15 * time.Minute),
	)
	codeVerifierStorage := ttlcache.New(
		ttlcache.WithTTL[string, string](15 * time.Minute),
	)
	go stateStorage.Start()
	go codeVerifierStorage.Start()

	return &Service{
		config:   oauth2Config,
		verifier: verifier,

		stateStorage:        stateStorage,
		codeVerifierStorage: codeVerifierStorage,
	}, nil
}

func (s *Service) AuthCodeURL(sessionID string, state string) string {
	s.stateStorage.Set(sessionID, state, ttlcache.DefaultTTL)

	codeVerifier := oauth2.GenerateVerifier()
	s.codeVerifierStorage.Set(sessionID, codeVerifier, ttlcache.DefaultTTL)

	return s.config.AuthCodeURL(state, oauth2.S256ChallengeOption(codeVerifier))
}

func (s *Service) VerifyState(sessionID, state string) bool {
	storedState, found := s.stateStorage.GetAndDelete(sessionID)
	if !found {
		return false
	}
	return storedState.Value() == state
}

func (s *Service) Exchange(ctx context.Context, sessionID, code string) (*oauth2.Token, error) {
	codeVerifier, found := s.codeVerifierStorage.GetAndDelete(sessionID)
	if !found {
		return nil, errors.New("code verifier not found")
	}
	token, err := s.config.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier.Value()))
	if err != nil {
		return nil, fmt.Errorf("exchange token: %w", err)
	}

	return token, nil
}

type UserInfo struct {
	ID   uuid.UUID
	Name string
}

func (s *Service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*UserInfo, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing id_token")
	}
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("verify id_token: %w", err)
	}

	type Payload struct {
		Sub  string `json:"sub"`
		Name string `json:"name"`
	}

	// traQ returns username in "name" claim
	// ref: https://github.com/traPtitech/traQ/blob/v3.21.0/service/oidc/userinfo.go#L57
	var payload Payload
	if err := idToken.Claims(&payload); err != nil {
		return nil, fmt.Errorf("parse id_token claims: %w", err)
	}

	id, err := uuid.Parse(payload.Sub)
	if err != nil {
		return nil, fmt.Errorf("parse user uuid: %w", err)
	}

	return &UserInfo{ID: id, Name: payload.Name}, nil
}
