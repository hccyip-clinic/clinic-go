package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/argon2"
)

// Argon2id parameters: 64 MiB memory, one iteration, four lanes, 32-byte key, 16-byte salt.
const (
	argonTime       = 1
	argonMemory     = 64 * 1024
	argonThreads    = 4
	argonKeyLen     = 32
	argonSaltLen    = 16
	sessionLifetime = 8 * time.Hour
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(key)), nil
}

func CheckPassword(encoded, password string) bool {
	var memory, iterations, threads uint32
	var saltText, keyText string
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false
	}
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &threads); err != nil {
		return false
	}
	saltText, keyText = parts[4], parts[5]
	salt, err1 := base64.RawStdEncoding.DecodeString(saltText)
	expected, err2 := base64.RawStdEncoding.DecodeString(keyText)
	if err1 != nil || err2 != nil || len(expected) == 0 {
		return false
	}
	actual := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(threads), uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

type Session struct {
	ID           string
	CSRFToken    string
	PasswordHash string
	CreatedAt    time.Time
}

func (s *Session) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:     "clinic_session",
		Value:    s.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(sessionLifetime / time.Second),
		Expires:  s.CreatedAt.Add(sessionLifetime),
	}
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func NewSessionStore() *SessionStore { return &SessionStore{sessions: make(map[string]Session)} }

func (s *SessionStore) Create(passwordHash ...string) Session {
	hash := ""
	if len(passwordHash) > 0 {
		hash = passwordHash[0]
	}
	session := Session{ID: randomToken(), CSRFToken: randomToken(), PasswordHash: hash, CreatedAt: time.Now()}
	s.mu.Lock()
	s.sessions[session.ID] = session
	s.mu.Unlock()
	return session
}

func (s *SessionStore) Get(id string) (Session, bool) {
	s.mu.RLock()
	session, ok := s.sessions[id]
	s.mu.RUnlock()
	if ok && time.Since(session.CreatedAt) >= sessionLifetime {
		s.Delete(id)
		return Session{}, false
	}
	return session, ok
}

func (s *SessionStore) Delete(id string) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

func (s *SessionStore) InvalidateAll() {
	s.mu.Lock()
	s.sessions = make(map[string]Session)
	s.mu.Unlock()
}

func ValidateCSRF(req *http.Request, session Session) bool {
	token := req.Header.Get("X-CSRF-Token")
	if token == "" {
		token = req.FormValue("csrf_token")
	}
	return token != "" && subtle.ConstantTimeCompare([]byte(token), []byte(session.CSRFToken)) == 1
}

func randomToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
