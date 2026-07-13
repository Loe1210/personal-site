package xauth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"

	"github.com/Loe1210/personal-site/configs"
)

const (
	placeholderBackend = "monolith-memory"
	claimsContextKey   = "session_claims"
)

var errSessionNotFound = errors.New("session not found")

var (
	sessionMu    sync.RWMutex
	sessionStore = make(map[string]Claims)
)

type SessionMetadata struct {
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	TraceID   string    `json:"trace_id,omitempty"`
}

type Claims struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	SessionMetadata
}

type SessionBundle struct {
	SessionID  string    `json:"session_id"`
	CookieName string    `json:"cookie_name"`
	Backend    string    `json:"backend"`
	ExpiresAt  time.Time `json:"expires_at"`
	TraceID    string    `json:"trace_id,omitempty"`
}

type storeConfig struct {
	Prefix     string
	ExpireHour int
	CookieName string
}

func CreateSession(userID int64, username string, roles []string) (string, error) {
	bundle, err := CreateSessionBundle(userID, username, roles)
	if err != nil {
		return "", err
	}
	return bundle.SessionID, nil
}

func CreateSessionBundle(userID int64, username string, roles []string) (*SessionBundle, error) {
	return CreateSessionBundleWithTrace(userID, username, roles, "")
}

func CreateSessionBundleWithTrace(userID int64, username string, roles []string, traceID string) (*SessionBundle, error) {
	cfg, err := loadStoreConfig()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    append([]string(nil), roles...),
		SessionMetadata: SessionMetadata{
			CreatedAt: now,
			ExpiresAt: now.Add(time.Duration(cfg.ExpireHour) * time.Hour),
			TraceID:   traceID,
		},
	}

	sessionID, err := newSessionID(cfg.Prefix)
	if err != nil {
		return nil, err
	}

	sessionMu.Lock()
	sessionStore[sessionID] = claims
	sessionMu.Unlock()

	return &SessionBundle{
		SessionID:  sessionID,
		CookieName: cfg.CookieName,
		Backend:    placeholderBackend,
		ExpiresAt:  claims.ExpiresAt,
		TraceID:    traceID,
	}, nil
}

func newSessionID(prefix string) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return prefix + hex.EncodeToString(randomBytes), nil
}
func ParseSession(raw string) (*Claims, error) {
	if raw == "" {
		return nil, errors.New("empty session id")
	}

	sessionMu.RLock()
	claims, ok := sessionStore[raw]
	sessionMu.RUnlock()
	if !ok {
		return nil, errSessionNotFound
	}

	if time.Now().After(claims.ExpiresAt) {
		_ = DestroySession(raw)
		return nil, errSessionNotFound
	}

	copyClaims := claims
	copyClaims.Roles = append([]string(nil), claims.Roles...)
	return &copyClaims, nil
}

func DestroySession(raw string) error {
	if raw == "" {
		return nil
	}

	sessionMu.Lock()
	delete(sessionStore, raw)
	sessionMu.Unlock()
	return nil
}

func SessionIDFromRequest(c *app.RequestContext) string {
	return string(c.Cookie(SessionCookieName()))
}

func SessionCookieName() string {
	cfg, err := loadStoreConfig()
	if err != nil {
		return "session_id"
	}
	return cfg.CookieName
}

func WriteSessionCookie(c *app.RequestContext, bundle *SessionBundle) {
	maxAge := int(time.Until(bundle.ExpiresAt).Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	c.SetCookie(bundle.CookieName, bundle.SessionID, maxAge, "/", "", protocol.CookieSameSiteLaxMode, false, true)
}

func ClearSessionCookie(c *app.RequestContext) {
	c.SetCookie(SessionCookieName(), "", -1, "/", "", protocol.CookieSameSiteLaxMode, false, true)
}

func SetClaims(c *app.RequestContext, claims *Claims) {
	c.Set(claimsContextKey, claims)
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("roles", append([]string(nil), claims.Roles...))
}

func ClaimsFromContext(c *app.RequestContext) (*Claims, bool) {
	value, ok := c.Get(claimsContextKey)
	if !ok {
		return nil, false
	}
	claims, ok := value.(*Claims)
	return claims, ok
}

func loadStoreConfig() (*storeConfig, error) {
	cfg := configs.AppConfig
	if cfg == nil {
		loaded, err := configs.Load("")
		if err != nil {
			return nil, err
		}
		cfg = loaded
	}

	return &storeConfig{
		Prefix:     cfg.SessionStore.Prefix,
		ExpireHour: cfg.SessionStore.ExpireHour,
		CookieName: cfg.SessionStore.CookieName,
	}, nil
}
