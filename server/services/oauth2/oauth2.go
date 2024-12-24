package oauth2

import (
	"context"
	"errors"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
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
		return nil, err
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	stateStorage := ttlcache.New[string, string]()
	codeVerifierStorage := ttlcache.New[string, string]()
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
	s.stateStorage.Set(sessionID, state, 15*time.Minute)

	codeVerifier := oauth2.GenerateVerifier()
	s.codeVerifierStorage.Set(sessionID, codeVerifier, 15*time.Minute)

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
	return s.config.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier.Value()))
}

type TraQUserInfo struct {
	ID   string `json:"sub"`
	Name string `json:"name"`
}

func (s *Service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*TraQUserInfo, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing id_token")
	}
	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	// traQ returns username in "name" claim
	// ref: https://github.com/traPtitech/traQ/blob/v3.21.0/service/oidc/userinfo.go#L57
	var info TraQUserInfo
	if err := idToken.Claims(&info); err != nil {
		return nil, err
	}

	return &info, nil
}
