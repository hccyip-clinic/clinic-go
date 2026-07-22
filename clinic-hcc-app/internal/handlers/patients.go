package handlers

import (
	"net/http"

	"clinic-hcc-app/internal/models"
)

func (r *Router) PatientList(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("q")
	searchQuery := query
	if canonical, err := models.NormalizeHKID(query); err == nil {
		searchQuery = canonical
	}
	patients, err := r.patients.Search(req.Context(), searchQuery)
	if err != nil {
		http.Error(w, "unable to load patients", http.StatusInternalServerError)
		return
	}
	r.render(w, "patients", map[string]interface{}{
		"Title":      "Patients",
		"ActivePage": "patients",
		"Patients":   patients,
		"Query":      query,
	})
}

func (r *Router) PatientFormNew(w http.ResponseWriter, req *http.Request) {
	r.render(w, "patient-form", map[string]interface{}{
		"Title":      "New Patient",
		"ActivePage": "patients",
		"Patient":    &models.Patient{Gender: "O"},
	})
}

func (r *Router) PatientCreate(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	patient := &models.Patient{
		Name:   req.FormValue("name"),
		HKID:   req.FormValue("hkid"),
		Gender: req.FormValue("gender"),
	}
	if err := r.patients.Create(req.Context(), patient); err != nil {
		r.render(w, "patient-form", map[string]interface{}{
			"Title":      "New Patient",
			"ActivePage": "patients",
			"Patient":    patient,
			"Error":      err.Error(),
		})
		return
	}
	http.Redirect(w, req, "/patients", http.StatusSeeOther)
}
