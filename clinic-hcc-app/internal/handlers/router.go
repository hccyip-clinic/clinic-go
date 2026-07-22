package handlers

import (
	"net/http"
	"time"

	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/repository"
	"clinic-hcc-app/internal/security"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	db       *database.DB
	patients *repository.PatientRepository
	receipts *repository.ReceiptRepository
	sessions *security.SessionStore
}

func NewRouter(db *database.DB) *Router {
	return &Router{
		db:       db,
		patients: repository.NewPatientRepository(db),
		receipts: repository.NewReceiptRepository(db),
		sessions: security.NewSessionStore(),
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

	mux.Get("/setup", r.SetupForm)
	mux.Post("/setup", r.SetupPassword)
	mux.Get("/login", r.LoginForm)
	mux.Post("/login", r.Login)
	mux.Group(func(protected chi.Router) {
		protected.Use(r.requireAuth)
		protected.Post("/logout", r.Logout)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsRead, next) }).Get("/", r.Dashboard)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsRead, next) }).Get("/receipts", r.ReceiptList)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsCreate, next) }).Get("/receipts/new", r.ReceiptFormNew)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsCreate, next) }).Post("/receipts", r.ReceiptCreate)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsRead, next) }).Get("/receipts/{id}", r.ReceiptView)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsUpdate, next) }).Get("/receipts/{id}/edit", r.ReceiptFormEdit)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsUpdate, next) }).Post("/receipts/{id}", r.ReceiptUpdate)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsArchive, next) }).Delete("/receipts/{id}", r.ReceiptDelete)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsArchive, next) }).Post("/receipts/{id}/delete", r.ReceiptDelete)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermReceiptsFinalize, next) }).Post("/receipts/{id}/finalize", r.ReceiptFinalize)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsRead, next) }).Get("/patients", r.PatientList)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsCreate, next) }).Get("/patients/new", r.PatientFormNew)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsCreate, next) }).Post("/patients", r.PatientCreate)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsUpdate, next) }).Get("/patients/{id}/edit", r.PatientFormEdit)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsUpdate, next) }).Post("/patients/{id}", r.PatientUpdate)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermPatientsRead, next) }).Get("/patients/search", r.PatientList)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermSettingsRead, next) }).Get("/settings", r.Settings)
		protected.With(func(next http.Handler) http.Handler { return r.permission(PermSettingsUpdate, next) }).Post("/settings", r.SettingsUpdate)
	})

	return mux
}
