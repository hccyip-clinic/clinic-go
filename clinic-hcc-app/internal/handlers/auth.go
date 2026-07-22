package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"clinic-hcc-app/internal/models"
	"clinic-hcc-app/internal/security"
)

func (r *Router) SetupForm(w http.ResponseWriter, req *http.Request) {
	r.renderPublic(w, "setup", map[string]interface{}{"Title": "Create password"})
}

func (r *Router) SetupPassword(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if hash, _ := r.passwordHash(); hash != "" {
		http.Error(w, "setup already completed", http.StatusConflict)
		return
	}
	password := req.FormValue("password")
	if len(password) < 12 || password != req.FormValue("confirm_password") {
		http.Error(w, "password must be at least 12 characters and match confirmation", http.StatusBadRequest)
		return
	}
	hash, err := security.HashPassword(password)
	if err != nil {
		http.Error(w, "unable to create password", http.StatusInternalServerError)
		return
	}
	if _, err := r.db.ExecContext(req.Context(), `UPDATE settings SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1`, hash); err != nil {
		http.Error(w, "unable to save password", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (r *Router) LoginForm(w http.ResponseWriter, req *http.Request) {
	r.renderPublic(w, "login", map[string]interface{}{"Title": "Sign in"})
}

func (r *Router) Login(w http.ResponseWriter, req *http.Request) {
	hash, err := r.passwordHash()
	if err != nil || hash == "" || !security.CheckPassword(hash, req.FormValue("password")) {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}
	session := r.sessions.Create(hash)
	http.SetCookie(w, session.Cookie())
	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Router) Logout(w http.ResponseWriter, req *http.Request) {
	if cookie, err := req.Cookie("clinic_session"); err == nil {
		r.sessions.Delete(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "clinic_session", Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteStrictMode})
	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (r *Router) passwordHash() (string, error) {
	var hash string
	err := r.db.QueryRow(`SELECT password_hash FROM settings WHERE id = 1`).Scan(&hash)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return hash, err
}

func (r *Router) renderPublic(w http.ResponseWriter, page string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles("templates/" + page + ".html")
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (r *Router) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("clinic_session")
		if err != nil {
			if hash, _ := r.passwordHash(); hash == "" {
				http.Redirect(w, req, "/setup", http.StatusSeeOther)
			} else {
				http.Redirect(w, req, "/login", http.StatusSeeOther)
			}
			return
		}
		session, ok := r.sessions.Get(cookie.Value)
		if !ok {
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		if hash, err := r.passwordHash(); err != nil || hash == "" || session.PasswordHash != hash {
			r.sessions.InvalidateAll()
			http.Redirect(w, req, "/login", http.StatusSeeOther)
			return
		}
		if req.Method != http.MethodGet && req.Method != http.MethodHead && req.Method != http.MethodOptions && !security.ValidateCSRF(req, session) {
			http.Error(w, "csrf validation failed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (r *Router) permission(perm models.Permission, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !models.DefaultUser().HasPermission(perm) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, req)
	})
}
