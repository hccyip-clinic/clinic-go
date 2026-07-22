package handlers

import (
	"net/http"
	"time"

	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	db       *database.DB
	patients *repository.PatientRepository
	receipts *repository.ReceiptRepository
}

func NewRouter(db *database.DB) *Router {
	return &Router{
		db:       db,
		patients: repository.NewPatientRepository(db),
		receipts: repository.NewReceiptRepository(db),
	}
}

func (r *Router) Setup() http.Handler {
	mux := chi.NewMux()

	// Middleware
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))

	// Static files
	mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Routes
	mux.Get("/", r.Dashboard)
	mux.Get("/receipts", r.ReceiptList)
	mux.Get("/receipts/new", r.ReceiptFormNew)
	mux.Post("/receipts", r.ReceiptCreate)
	mux.Get("/receipts/{id}", r.ReceiptView)
	mux.Get("/receipts/{id}/edit", r.ReceiptFormEdit)
	mux.Post("/receipts/{id}", r.ReceiptUpdate)
	mux.Delete("/receipts/{id}", r.ReceiptDelete)
	mux.Post("/receipts/{id}/delete", r.ReceiptDelete)
	mux.Post("/receipts/{id}/finalize", r.ReceiptFinalize)
	mux.Get("/patients", r.PatientList)
	mux.Get("/patients/new", r.PatientFormNew)
	mux.Post("/patients", r.PatientCreate)
	mux.Get("/patients/{id}/edit", r.PatientFormEdit)
	mux.Post("/patients/{id}", r.PatientUpdate)
	mux.Get("/patients/search", r.PatientList)
	mux.Get("/settings", r.Settings)
	mux.Post("/settings", r.SettingsUpdate)

	return mux
}
