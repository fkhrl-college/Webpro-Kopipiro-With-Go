package main

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

type Session struct {
	UserID   int
	Username string
	Role     string
}

func (s Session) IsAdmin() bool {
	return s.Role == "admin"
}

const cookieName = "session_id"

var (
	sessions  = map[string]Session{}
	sessionMu sync.RWMutex
)

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func startSession(w http.ResponseWriter, s Session) {
	token := newToken()
	sessionMu.Lock()
	sessions[token] = s
	sessionMu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func currentSession(r *http.Request) (Session, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return Session{}, false
	}

	sessionMu.RLock()
	defer sessionMu.RUnlock()
	s, ok := sessions[c.Value]
	return s, ok
}

func endSession(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(cookieName)
	if err == nil {
		sessionMu.Lock()
		delete(sessions, c.Value)
		sessionMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1})
}

// middleware
type protectedHandler func(w http.ResponseWriter, r *http.Request, s Session)

func requireLogin(next protectedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := currentSession(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r, s)
	}
}

func requireAdmin(next protectedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := currentSession(r)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// 403 forbidden access
		if !s.IsAdmin() {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next(w, r, s)
	}
}
