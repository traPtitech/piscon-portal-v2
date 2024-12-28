package internal

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
)

type SessionManager struct {
	secret string
	debug  bool
}

const sessionName = "picon-portal-v2-session"
const sessionIDKey = "session_id"

func NewSessionManager(secret string, debug bool) *SessionManager {
	return &SessionManager{
		secret: secret,
		debug:  debug,
	}
}

func (sm *SessionManager) Init(e *echo.Group) {
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(sm.secret))))
}

func getSession(c echo.Context) (*sessions.Session, error) {
	return session.Get(sessionName, c)
}

func (sm *SessionManager) GetSessionID(c echo.Context) (string, error) {
	sess, err := getSession(c)
	if err != nil {
		return "", err
	}
	sessID, ok := sess.Values[sessionIDKey].(string)
	if !ok {
		return "", nil
	}
	return sessID, nil
}

func (sm *SessionManager) SetSessionID(c echo.Context, maxAge time.Duration) (string, error) {
	sess, err := getSession(c)
	if err != nil {
		return "", err
	}
	sessID := domain.NewSessionID()
	sess.Values[sessionIDKey] = sessID
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

func (sm *SessionManager) ClearSessionID(c echo.Context) error {
	sess, err := getSession(c)
	if err != nil {
		return err
	}
	delete(sess.Values, sessionIDKey)
	return sess.Save(c.Request(), c.Response())
}
