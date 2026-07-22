package security

import (
	"net/http/httptest"
	"testing"
)

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "correct horse battery staple") {
		t.Fatal("expected password to verify")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("wrong password verified")
	}
}

func TestSessionStoreInvalidatesSessions(t *testing.T) {
	store := NewSessionStore()
	session := store.Create()
	if _, ok := store.Get(session.ID); !ok {
		t.Fatal("session was not stored")
	}

	store.InvalidateAll()
	if _, ok := store.Get(session.ID); ok {
		t.Fatal("session survived invalidation")
	}
}

func TestSessionCookieExpires(t *testing.T) {
	session := NewSessionStore().Create()
	cookie := session.Cookie()
	if cookie.MaxAge <= 0 || cookie.Expires.IsZero() {
		t.Fatal("expected session cookie expiry")
	}
}

func TestCSRFTokenIsRequiredForStateChanges(t *testing.T) {
	store := NewSessionStore()
	session := store.Create()
	req := httptest.NewRequest("POST", "/", nil)
	req.AddCookie(session.Cookie())
	if ValidateCSRF(req, session) {
		t.Fatal("request without csrf token was accepted")
	}
	req.Header.Set("X-CSRF-Token", session.CSRFToken)
	if !ValidateCSRF(req, session) {
		t.Fatal("request with csrf token was rejected")
	}
}
