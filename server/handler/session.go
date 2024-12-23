package handler

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type sessionManager struct {
	secret string
	debug  bool
}

const sessionName = "picon-portal-v2-session"

func newSessionManager(secret string, debug bool) *sessionManager {
	return &sessionManager{
		secret: secret,
		debug:  debug,
	}
}

func (sm *sessionManager) init(e *echo.Group) {
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sm.secret))))
}

func getSession(c echo.Context) (*sessions.Session, error) {
	return session.Get(sessionName, c)
}

func (sm *sessionManager) getSessionID(c echo.Context) (string, error) {
	sess, err := getSession(c)
	if err != nil {
		return "", err
	}
	sessID, ok := sess.Values["session_id"].(string)
	if !ok {
		return "", nil
	}
	return sessID, nil
}

func (sm *sessionManager) setSessionID(c echo.Context, maxAge time.Duration) (string, error) {
	sess, err := getSession(c)
	if err != nil {
		return "", err
	}
	sessID := domain.NewSessionID()
	sess.Values["session_id"] = sessID
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int(maxAge.Milliseconds() / 1000),
		HttpOnly: true,
		Secure:   !sm.debug,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return "", err
	}
	return sessID, nil
}

func (sm *sessionManager) clearSessionID(c echo.Context) error {
	sess, err := getSession(c)
	if err != nil {
		return err
	}
	delete(sess.Values, "session_id")
	return sess.Save(c.Request(), c.Response())
}
