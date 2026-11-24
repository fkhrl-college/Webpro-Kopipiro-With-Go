package main

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginPage struct {
	Error string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	_, ok := currentSession(r)
	if ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	page := LoginPage{}
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		var s Session
		var hash string
		err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = ?", username).Scan(&s.UserID, &s.Username, &hash, &s.Role)
		if err == nil && bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil {
			startSession(w, s)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if err != nil && err != sql.ErrNoRows {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		page.Error = "Incorrect username or password."
	}
	tpl.ExecuteTemplate(w, "login.html", page)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	endSession(w, r)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
