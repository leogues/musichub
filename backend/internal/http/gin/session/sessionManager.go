package session

import (
	"encoding/hex"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

const SessionCookieName = "session"

type HashedBlock struct {
	HashKey  string
	BlockKey string
}

type SessionManager struct {
	sc *securecookie.SecureCookie
}

func NewSessionManager(hashKey string, blockKey string) (*SessionManager, error) {
	sc, err := openSecureCookie(HashedBlock{hashKey, blockKey})

	if err != nil {
		return nil, err
	}

	return &SessionManager{sc}, nil
}

func openSecureCookie(h HashedBlock) (*securecookie.SecureCookie, error) {
	if h.HashKey == "" {
		return nil, errors.New("HashKey is required")
	}

	if h.BlockKey == "" {
		return nil, errors.New("BlockKey is required")
	}

	hashKey, err := hex.DecodeString(h.HashKey)
	if err != nil {
		return nil, errors.New("invalid hash key")
	}

	blockKey, err := hex.DecodeString(h.BlockKey)
	if err != nil {
		return nil, errors.New("invalid block key")
	}

	sc := securecookie.New(hashKey, blockKey)
	sc.SetSerializer(securecookie.JSONEncoder{})

	return sc, nil
}

// Session returns the session from the cookie
func (s *SessionManager) Session(c *gin.Context) (Session, error) {

	cookie, err := c.Cookie(SessionCookieName)

	if err != nil {
		return NewSession(), nil
	}

	var session Session
	if err := s.UnmarshalSession(cookie, &session); err != nil {

		return NewSession(), err
	}

	return session, nil
}

// SetSession sets the session in the cookie
func (s *SessionManager) SetSession(c *gin.Context, session Session) error {
	buf, err := s.MarshalSession(session)

	if err != nil {
		return err
	}

	expires := 3600 * 24 * 30 // 1 hour * 24 hours * 30 days

	c.SetCookie(SessionCookieName, buf, expires, "/", "", false, true)
	return nil
}

func (s *SessionManager) MarshalSession(session Session) (string, error) {
	return s.sc.Encode(SessionCookieName, session)
}

func (s *SessionManager) UnmarshalSession(data string, session *Session) error {
	return s.sc.Decode(SessionCookieName, data, session)
}
