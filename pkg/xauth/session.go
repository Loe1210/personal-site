package xauth

import (
	"context"
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
	memoryBackend    = "monolith-memory"
	redisBackend     = "redis"
	claimsContextKey = "session_claims"
)

var errSessionNotFound = errors.New("session not found")

type Store interface {
	Save(ctx context.Context, sessionID string, claims *Claims, ttl time.Duration) error
	Get(ctx context.Context, sessionID string) (*Claims, error)
	Delete(ctx context.Context, sessionID string) error
	Backend() string
}

var activeStore Store = newMemoryStoreForTest()

type memoryStore struct {
	mu       sync.RWMutex
	sessions map[string]Claims
}

func newMemoryStoreForTest() Store {
	return &memoryStore{sessions: make(map[string]Claims)}
}

func (s *memoryStore) Save(_ context.Context, sessionID string, claims *Claims, _ time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyClaims := *claims
	copyClaims.Roles = append([]string(nil), claims.Roles...)
	s.sessions[sessionID] = copyClaims
	return nil
}

func (s *memoryStore) Get(_ context.Context, sessionID string) (*Claims, error) {
	s.mu.RLock()
	claims, ok := s.sessions[sessionID]
	s.mu.RUnlock()
	if !ok {
		return nil, errSessionNotFound
	}
	if time.Now().After(claims.ExpiresAt) {
		_ = s.Delete(context.Background(), sessionID)
		return nil, errSessionNotFound
	}
	copyClaims := claims
	copyClaims.Roles = append([]string(nil), claims.Roles...)
	return &copyClaims, nil
}

func (s *memoryStore) Delete(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
	return nil
}

func (s *memoryStore) Backend() string {
	return memoryBackend
}

func UseStore(store Store) {
	if store == nil {
		activeStore = newMemoryStoreForTest()
		return
	}
	activeStore = store
}

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
	return CreateSessionWithContext(context.Background(), userID, username, roles)
}

func CreateSessionWithContext(ctx context.Context, userID int64, username string, roles []string) (string, error) {
	bundle, err := CreateSessionBundleWithContext(ctx, userID, username, roles)
	if err != nil {
		return "", err
	}
	return bundle.SessionID, nil
}

func CreateSessionBundle(userID int64, username string, roles []string) (*SessionBundle, error) {
	return CreateSessionBundleWithContext(context.Background(), userID, username, roles)
}

func CreateSessionBundleWithContext(ctx context.Context, userID int64, username string, roles []string) (*SessionBundle, error) {
	return CreateSessionBundleWithTraceContext(ctx, userID, username, roles, "")
}

func CreateSessionBundleWithTrace(userID int64, username string, roles []string, traceID string) (*SessionBundle, error) {
	return CreateSessionBundleWithTraceContext(context.Background(), userID, username, roles, traceID)
}

func CreateSessionBundleWithTraceContext(ctx context.Context, userID int64, username string, roles []string, traceID string) (*SessionBundle, error) {
	cfg, err := loadStoreConfig()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	ttl := time.Duration(cfg.ExpireHour) * time.Hour
	claims := Claims{
		UserID:   userID,
		Username: username,
		Roles:    append([]string(nil), roles...),
		SessionMetadata: SessionMetadata{
			CreatedAt: now,
			ExpiresAt: now.Add(ttl),
			TraceID:   traceID,
		},
	}

	sessionID, err := newSessionID(cfg.Prefix)
	if err != nil {
		return nil, err
	}
	if err := activeStore.Save(ctx, sessionID, &claims, ttl); err != nil {
		return nil, err
	}

	return &SessionBundle{
		SessionID:  sessionID,
		CookieName: cfg.CookieName,
		Backend:    activeStore.Backend(),
		ExpiresAt:  claims.ExpiresAt,
		TraceID:    traceID,
	}, nil
}

func newSessionID(_ string) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(randomBytes), nil
}

func ParseSession(raw string) (*Claims, error) {
	return ParseSessionWithContext(context.Background(), raw)
}

func ParseSessionWithContext(ctx context.Context, raw string) (*Claims, error) {
	if raw == "" {
		return nil, errors.New("empty session id")
	}
	return activeStore.Get(ctx, raw)
}

func DestroySession(raw string) error {
	return DestroySessionWithContext(context.Background(), raw)
}

func DestroySessionWithContext(ctx context.Context, raw string) error {
	if raw == "" {
		return nil
	}
	return activeStore.Delete(ctx, raw)
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
