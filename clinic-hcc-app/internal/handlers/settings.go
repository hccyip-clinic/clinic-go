package handlers

import "net/http"

import "clinic-hcc-app/internal/security"

func (r *Router) Settings(w http.ResponseWriter, req *http.Request) {
	var settings struct {
		ClinicName   string
		Address      string
		Phone        string
		Practitioner string
	}
	err := r.db.QueryRowContext(req.Context(), `
		SELECT clinic_name, clinic_address, clinic_phone, practitioner_name
		FROM settings WHERE id = 1
	`).Scan(&settings.ClinicName, &settings.Address, &settings.Phone, &settings.Practitioner)
	if err != nil {
		http.Error(w, "unable to load settings", http.StatusInternalServerError)
		return
	}
	r.render(w, req, "settings", map[string]interface{}{
		"Title":      "Settings",
		"ActivePage": "settings",
		"Settings":   settings,
	})
}

func (r *Router) SettingsUpdate(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	newPassword := req.FormValue("new_password")
	if newPassword != "" {
		if len(newPassword) < 12 || newPassword != req.FormValue("confirm_password") {
			http.Error(w, "password must be at least 12 characters and match confirmation", http.StatusBadRequest)
			return
		}
		hash, err := security.HashPassword(newPassword)
		if err != nil {
			http.Error(w, "unable to create password", http.StatusInternalServerError)
			return
		}
		if _, err := r.db.ExecContext(req.Context(), `UPDATE settings SET password_hash = ? WHERE id = 1`, hash); err != nil {
			http.Error(w, "unable to update password", http.StatusInternalServerError)
			return
		}
		r.sessions.InvalidateAll()
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	_, err := r.db.ExecContext(req.Context(), `
		UPDATE settings SET clinic_name = ?, clinic_address = ?, clinic_phone = ?,
			practitioner_name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1
	`, req.FormValue("clinic_name"), req.FormValue("address"), req.FormValue("phone"),
		req.FormValue("practitioner"))
	if err != nil {
		http.Error(w, "unable to update settings", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, req, "/settings", http.StatusSeeOther)
}
