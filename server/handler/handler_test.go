package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
	"github.com/traPtitech/piscon-portal-v2/server/repository"
	"github.com/traPtitech/piscon-portal-v2/server/services/oauth2"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"github.com/traPtitech/piscon-portal-v2/server/utils/random"
)

var oauth2ServerURL string

func TestMain(m *testing.M) {
	oauth2Server := newOauth2MockServer()
	defer oauth2Server.Close()

	m.Run()
}

func NewPortalServer(useCase usecase.UseCase, repo repository.Repository) *httptest.Server {
	e := echo.New()
	server := httptest.NewTLSServer(e)

	config := handler.Config{
		RootURL:       server.URL,
		SessionSecret: "secret",
		Oauth2: oauth2.Config{
			Issuer:       oauth2ServerURL,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AuthURL:      oauth2ServerURL + "/authorize",
			TokenURL:     oauth2ServerURL + "/token",
		},
	}
	h, err := handler.New(useCase, repo, config)
	if err != nil {
		panic(err)
	}

	h.SetupRoutes(e)

	return server
}

func NewHandler(useCase usecase.UseCase, repo repository.Repository, sessionManager handler.SessionManager) *handler.Handler {
	config := handler.Config{
		RootURL:       "http://localhost",
		SessionSecret: "secret",
		Oauth2: oauth2.Config{
			Issuer:       oauth2ServerURL,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AuthURL:      oauth2ServerURL + "/authorize",
			TokenURL:     oauth2ServerURL + "/token",
		},
		SessionManager: sessionManager,
	}
	h, err := handler.New(useCase, repo, config)
	if err != nil {
		panic(err)
	}

	return h
}

func NewClient(server *httptest.Server) *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: server.Client().Transport,
		Jar:       jar,
	}
	return client
}

func Login(t *testing.T, server *httptest.Server, client *http.Client, userID uuid.UUID) error {
	t.Helper()

	// not following redirect for the first request
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	codeRes, err := client.Get(joinPath(t, server.URL, "/api/oauth2/code"))
	if err != nil {
		t.Error(err)
		return fmt.Errorf("get authorization URL: %w", err)
	}
	defer codeRes.Body.Close()
	if codeRes.StatusCode != http.StatusSeeOther {
		msg := fmt.Sprintf("unexpected status code: expected %d, got %d", http.StatusSeeOther, codeRes.StatusCode)
		t.Error(msg)
		return errors.New(msg)
	}
	// set username for testing
	authURL, err := url.Parse(codeRes.Header.Get("Location"))
	if err != nil {
		t.Error(err)
		return fmt.Errorf("parse authorization URL: %w", err)
	}
	q := authURL.Query()
	q.Add("user", userID.String())
	authURL.RawQuery = q.Encode()

	// from here, follow redirect
	client.CheckRedirect = nil
	res, err := client.Get(authURL.String())
	if err != nil {
		t.Error(err)
		return fmt.Errorf("get authorization URL: %w", err)
	}
	defer res.Body.Close()
	// status code should be 200 or 404 (redirected to root page of frontend)
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNotFound {
		msg := fmt.Sprintf("unexpected status code: expected %d or %d, got %d", http.StatusOK, http.StatusNotFound, res.StatusCode)
		t.Error(msg)
		return errors.New(msg)
	}

	return nil
}

func joinPath(t *testing.T, base, path string) string {
	t.Helper()
	res, err := url.JoinPath(base, path)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

type oauth2MockServer struct {
	key     *rsa.PrivateKey
	rootURL string

	codeChallengeMap *ttlcache.Cache[string, string]
	userMap          *ttlcache.Cache[string, string]
}

func newOauth2MockServer() *httptest.Server {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)
	oauth2ServerURL = server.URL

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	s := &oauth2MockServer{
		key:              key,
		rootURL:          oauth2ServerURL,
		codeChallengeMap: ttlcache.New[string, string](),
		userMap:          ttlcache.New[string, string](),
	}

	mux.HandleFunc("/.well-known/openid-configuration", s.handleWellKnown)
	mux.HandleFunc("/jwks", s.handleJWKS)
	mux.HandleFunc("/authorize", s.handleAuthorize)
	mux.HandleFunc("/token", s.handleToken)

	return server
}

func (s *oauth2MockServer) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	redirectURI := q.Get("redirect_uri")
	codeChallenge := q.Get("code_challenge")
	state := q.Get("state")

	code := random.String(32)

	s.codeChallengeMap.Set(code, codeChallenge, 15*time.Minute)
	// use user params as username for testing
	if user := q.Get("user"); user != "" {
		s.userMap.Set(code, user, 15*time.Minute)
	}

	http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state), http.StatusSeeOther)
}

func (s *oauth2MockServer) handleToken(w http.ResponseWriter, r *http.Request) {
	codeVerifier := r.FormValue("code_verifier")
	code := r.FormValue("code")

	hash := sha256.Sum256([]byte(codeVerifier))
	encoded := base64.RawURLEncoding.EncodeToString(hash[:])
	v := s.codeChallengeMap.Get(code)
	if v == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if v.Value() != encoded {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var userID string
	user, found := s.userMap.GetAndDelete(code)
	if !found {
		userID = uuid.NewString()
	} else {
		userID = user.Value()
	}

	claims := struct {
		jwt.RegisteredClaims
		Name string `json:"name"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.rootURL,
			Subject:   userID,
			Audience:  []string{"client-id"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Name: "test-user",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	idToken, err := token.SignedString(s.key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
		IDToken     string `json:"id_token"`
	}{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		Scope:       "openid",
		IDToken:     idToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (s *oauth2MockServer) handleWellKnown(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"issuer":                 s.rootURL,
		"authorization_endpoint": s.rootURL + "/authorize",
		"token_endpoint":         s.rootURL + "/token",
		"jwks_uri":               s.rootURL + "/jwks",
	})
}

func (s *oauth2MockServer) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"keys": []map[string]any{
			{
				"kty": "RSA",
				"alg": "RS256",
				"use": "sig",
				"kid": "key-id",
				"n":   base64.RawURLEncoding.EncodeToString(s.key.N.Bytes()),
				"e":   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(s.key.E)).Bytes()),
			},
		},
	})
}
