package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"clinic-hcc-app/internal/models"
)

func (r *Router) render(w http.ResponseWriter, page string, data map[string]interface{}) {
	var clinicName string
	if err := r.db.QueryRow(`SELECT clinic_name FROM settings WHERE id = 1`).Scan(&clinicName); err != nil {
		http.Error(w, "unable to load clinic settings", http.StatusInternalServerError)
		return
	}
	data["ClinicName"] = clinicName
	data["User"] = models.DefaultUser()
	tmpl, err := template.New("layout.html").Funcs(template.FuncMap{
		"money": models.FormatMoney,
	}).ParseFiles(filepath.Join("templates", "layout.html"), filepath.Join("templates", page+".html"))
	if err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
