package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"clinic-hcc-app/internal/models"
)

func (r *Router) render(w http.ResponseWriter, page string, data map[string]interface{}) {
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
