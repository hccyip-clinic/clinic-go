package handlers

import (
	"fmt"
	"net/http"
	"time"

	"clinic-hcc-app/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	db *database.DB
}

func NewRouter(db *database.DB) *Router {
	return &Router{db: db}
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
	mux.Get("/patients", r.PatientList)
	mux.Get("/settings", r.Settings)
	mux.Post("/settings", r.SettingsUpdate)

	return mux
}

func (r *Router) Dashboard(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Dashboard - TODO")
}

func (r *Router) ReceiptList(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Receipt List - TODO")
}

func (r *Router) ReceiptFormNew(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "New Receipt Form - TODO")
}

func (r *Router) ReceiptCreate(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Create Receipt - TODO")
}

func (r *Router) ReceiptView(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "View Receipt - TODO", id)
}

func (r *Router) ReceiptFormEdit(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Edit Receipt Form - TODO", id)
}

func (r *Router) ReceiptUpdate(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Update Receipt - TODO", id)
}

func (r *Router) ReceiptDelete(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Delete Receipt - TODO", id)
}

func (r *Router) PatientList(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Patient List - TODO")
}

func (r *Router) Settings(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Settings - TODO")
}

func (r *Router) SettingsUpdate(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Update Settings - TODO")
}