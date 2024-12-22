package handler_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"database/sql"
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
	"github.com/jellydator/ttlcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/handler"
)

const Oauth2ServerURL = "http://localhost:9000"

var (
	ServerURL string

	db     *sql.DB
	server *httptest.Server
)

func TestMain(m *testing.M) {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/test?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	oauth2Server := newOauth2Server()
	defer oauth2Server.Close()

	server = newPortalServer()
	defer server.Close()

	m.Run()
}

func newPortalServer() *httptest.Server {
	e := echo.New()
	server := httptest.NewTLSServer(e)
	ServerURL = server.URL

	config := handler.Config{
		RootURL:       server.URL,
		SessionSecret: "secret",
		Oauth2: handler.Oauth2Config{
			Issuer:       Oauth2ServerURL,
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AuthURL:      Oauth2ServerURL + "/authorize",
			TokenURL:     Oauth2ServerURL + "/token",
		},
	}
	h, err := handler.New(db, config)
	if err != nil {
		panic(err)
	}

	h.SetupRoutes(e)

	return server
}

func NewClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: server.Client().Transport,
		Jar:       jar,
	}
	return client
}

func Login(t *testing.T, client *http.Client, userName string) error {
	t.Helper()

	// not following redirect for the first request
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	codeRes, err := client.Get(server.URL + "/api/oauth2/code")
	if err != nil {
		t.Error(err)
		return err
	}
	defer codeRes.Body.Close()
	if codeRes.StatusCode != http.StatusSeeOther {
		msg := fmt.Sprintf("unexpected status code: expected %d, got %d", http.StatusSeeOther, codeRes.StatusCode)
		t.Error(msg)
		return errors.New(msg)
	}
	// set username for testing
	authURL, _ := url.Parse(codeRes.Header.Get("Location"))
	q := authURL.Query()
	q.Add("user", userName)
	authURL.RawQuery = q.Encode()

	// from here, follow redirect
	client.CheckRedirect = nil
	res, err := client.Get(authURL.String())
	if err != nil {
		t.Error(err)
		return err
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

type oauth2Server struct {
	key     *rsa.PrivateKey
	rootURL string

	codeChallengeMap *ttlcache.Cache[string, string]
	userMap          *ttlcache.Cache[string, string]
}

func newOauth2Server() *http.Server {
	mux := http.NewServeMux()

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	s := &oauth2Server{
		key:              key,
		rootURL:          Oauth2ServerURL,
		codeChallengeMap: ttlcache.New[string, string](),
		userMap:          ttlcache.New[string, string](),
	}

	mux.HandleFunc("/.well-known/openid-configuration", s.handleWellKnown)
	mux.HandleFunc("/jwks", s.handleJWKS)
	mux.HandleFunc("/authorize", s.handleAuthorize)
	mux.HandleFunc("/token", s.handleToken)

	server := &http.Server{
		Addr:    ":9000",
		Handler: mux,
	}
	go server.ListenAndServe()

	return server
}

func (s *oauth2Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	redirectURI := q.Get("redirect_uri")
	codeChallenge := q.Get("code_challenge")
	state := q.Get("state")

	code := randomString()

	s.codeChallengeMap.Set(code, codeChallenge, 15*time.Minute)
	// use user params as username for testing
	if user := q.Get("user"); user != "" {
		s.userMap.Set(code, user, 15*time.Minute)
	}

	http.Redirect(w, r, fmt.Sprintf("%s?code=%s&state=%s", redirectURI, code, state), http.StatusSeeOther)
}

func (s *oauth2Server) handleToken(w http.ResponseWriter, r *http.Request) {
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

	var userName string
	user, found := s.userMap.GetAndDelete(code)
	if !found {
		userName = "test-user"
	} else {
		userName = user.Value()
	}

	claims := struct {
		jwt.RegisteredClaims
		Name string `json:"name"`
	}{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.rootURL,
			Subject:   userName,
			Audience:  []string{"client-id"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Name: userName,
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
	json.NewEncoder(w).Encode(res)
}

func (s *oauth2Server) handleWellKnown(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"issuer":                 s.rootURL,
		"authorization_endpoint": s.rootURL + "/authorize",
		"token_endpoint":         s.rootURL + "/token",
		"jwks_uri":               s.rootURL + "/jwks",
	})
}

func (s *oauth2Server) handleJWKS(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
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

func randomString() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
