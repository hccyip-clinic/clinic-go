# HTTP Server & TailwindCSS Setup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Set up the HTTP server with chi router, create handler scaffolding, configure TailwindCSS v4, and build the base dashboard layout.

**Architecture:** Server-rendered HTML using Go's html/template with HTMX for interactivity. TailwindCSS v4 via Vite plugin for styling. Chi router for HTTP routing.

**Tech Stack:** Go 1.26.5, chi v5, HTMX 2.x, TailwindCSS v4, @tailwindcss/vite plugin

## Global Constraints

- Go version: 1.26.5 (from go.mod)
- TailwindCSS v4 with CSS-first configuration (no tailwind.config.js)
- HTMX 2.x loaded from CDN
- All monetary values stored as integers (cents)
- SQLite with WAL mode enabled
- Server runs on localhost only (default port 8080)
- Templates use .html extension with Go template syntax
- Static files served from /static/ path

---

### Task 1: Add Go Dependencies (chi router)

**Files:**
- Modify: `clinic-hcc-app/go.mod`
- Modify: `clinic-hcc-app/go.sum` (auto-generated)

**Interfaces:**
- Consumes: Existing Go module structure
- Produces: chi router available for import in handlers

- [ ] **Step 1: Add chi router dependency**

Run in `clinic-hcc-app/` directory:
```bash
go get github.com/go-chi/chi/v5
```

Expected output:
```
go: added github.com/go-chi/chi/v5 v5.x.x
```

- [ ] **Step 2: Verify dependency added**

```bash
go mod tidy
```

Expected: No errors, go.sum updated

- [ ] **Step 3: Commit**

```bash
git add clinic-hcc-app/go.mod clinic-hcc-app/go.sum
git commit -m "chore: add chi router dependency"
```

---

### Task 2: Create HTTP Router and Middleware

**Files:**
- Modify: `clinic-hcc-app/cmd/server/main.go`
- Create: `clinic-hcc-app/internal/handlers/router.go`

**Interfaces:**
- Consumes: database.DB from Task 1
- Produces: chi router with middleware stack

- [ ] **Step 1: Create router.go with middleware**

```go
package handlers

import (
	"clinic-hcc-app/internal/database"
	"net/http"
	"time"

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

	// Routes (to be implemented in subsequent tasks)
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
```

- [ ] **Step 2: Add stub handlers to router.go**

```go
package handlers

import (
	"fmt"
	"net/http"
)

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
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "View Receipt - TODO")
}

func (r *Router) ReceiptFormEdit(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Edit Receipt Form - TODO")
}

func (r *Router) ReceiptUpdate(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Update Receipt - TODO")
}

func (r *Router) ReceiptDelete(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintln(w, "Delete Receipt - TODO")
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
```

- [ ] **Step 3: Update main.go to use router**

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"clinic-hcc-app/internal/config"
	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/handlers"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	router := handlers.NewRouter(db)
	mux := router.Setup()

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s (database: %s)", addr, cfg.DatabasePath)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```

- [ ] **Step 4: Build and verify compilation**

```bash
cd clinic-hcc-app
go build ./cmd/server
```

Expected: No errors, server.exe (or server binary) created

- [ ] **Step 5: Commit**

```bash
git add clinic-hcc-app/cmd/server/main.go clinic-hcc-app/internal/handlers/router.go
git commit -m "feat: set up chi router with middleware and stub handlers"
```

---

### Task 3: Set Up TailwindCSS v4 with Vite Plugin

**Files:**
- Create: `clinic-hcc-app/package.json`
- Create: `clinic-hcc-app/src/index.css`
- Modify: `clinic-hcc-app/static/css/.gitkeep` (ensure directory exists)

**Interfaces:**
- Consumes: Node.js 18+ runtime
- Produces: TailwindCSS build pipeline

- [ ] **Step 1: Create package.json**

```json
{
  "name": "clinic-hcc-app",
  "version": "1.0.0",
  "description": "Clinic receipt management system",
  "scripts": {
    "dev:css": "npx tailwindcss -i src/index.css -o static/css/styles.css --watch",
    "build:css": "npx tailwindcss -i src/index.css -o static/css/styles.css --minify"
  },
  "devDependencies": {
    "tailwindcss": "^4.0.0",
    "@tailwindcss/vite": "^4.0.0"
  }
}
```

- [ ] **Step 2: Create src/index.css with Tailwind v4 import**

```css
@import "tailwindcss";

@theme {
  --color-primary: #3b82f6;
  --color-secondary: #64748b;
  --color-accent: #10b981;
  --color-danger: #ef4444;
  --font-sans: 'Inter', system-ui, -apple-system, sans-serif;
  --radius-lg: 0.75rem;
  --radius-md: 0.375rem;
  --shadow-card: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1);
}

/* Custom utilities */
@utility sidebar-fixed {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: 16rem;
}

@utility content-auto {
  content-visibility: auto;
}

@utility card {
  background-color: white;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card);
  padding: 1.5rem;
}
```

- [ ] **Step 3: Install npm dependencies**

```bash
cd clinic-hcc-app
npm install
```

Expected output:
```
added 50 packages in 2m
```

- [ ] **Step 4: Build initial CSS**

```bash
npm run build:css
```

Expected: `static/css/styles.css` created with Tailwind utilities

- [ ] **Step 5: Commit**

```bash
git add clinic-hcc-app/package.json clinic-hcc-app/package-lock.json clinic-hcc-app/src/index.css clinic-hcc-app/static/css/styles.css
git commit -m "feat: add TailwindCSS v4 with CSS-first configuration"
```

---

### Task 4: Create Base Dashboard Layout Template

**Files:**
- Create: `clinic-hcc-app/templates/layouts/dashboard.html`
- Create: `clinic-hcc-app/templates/partials/sidebar.html`
- Create: `clinic-hcc-app/templates/partials/header.html`

**Interfaces:**
- Consumes: static/css/styles.css from Task 3
- Produces: Base layout for all pages

- [ ] **Step 1: Create dashboard.html layout**

```html
{{define "layouts/dashboard"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Clinic Management</title>
    <link rel="stylesheet" href="/static/css/styles.css">
    <script src="https://unpkg.com/htmx.org@2.0.0"></script>
</head>
<body class="bg-gray-100">
    <!-- Sidebar -->
    <aside class="sidebar-fixed bg-white shadow-lg z-10">
        {{template "partials/sidebar" .}}
    </aside>

    <!-- Main Content -->
    <div class="ml-64">
        <!-- Header -->
        <header class="bg-white shadow">
            {{template "partials/header" .}}
        </header>

        <!-- Page Content -->
        <main class="p-8">
            {{template "content" .}}
        </main>
    </div>

    <script>
        document.body.addEventListener('htmx:configRequest', (event) => {
            event.detail.headers['X-CSRF-Token'] = '{{.CSRFToken}}';
        });
    </script>
</body>
</html>
{{end}}
```

- [ ] **Step 2: Create sidebar.html partial**

```html
{{define "partials/sidebar"}}
<nav class="mt-6 px-4">
    <div class="mb-8">
        <h1 class="text-2xl font-bold text-gray-800">Clinic HCC</h1>
        <p class="text-sm text-gray-500">Management System</p>
    </div>

    <ul class="space-y-2">
        <li>
            <a href="/" class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 {{if eq .ActivePage "dashboard"}}bg-gray-100{{end}}">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                </svg>
                Dashboard
            </a>
        </li>
        <li>
            <a href="/receipts" class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 {{if eq .ActivePage "receipts"}}bg-gray-100{{end}}">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                </svg>
                Receipts
            </a>
        </li>
        <li>
            <a href="/patients" class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 {{if eq .ActivePage "patients"}}bg-gray-100{{end}}">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                Patients
            </a>
        </li>
        <li>
            <a href="/settings" class="flex items-center gap-3 px-4 py-3 rounded-lg text-gray-700 hover:bg-gray-100 {{if eq .ActivePage "settings"}}bg-gray-100{{end}}">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                </svg>
                Settings
            </a>
        </li>
    </ul>
</nav>
{{end}}
```

- [ ] **Step 3: Create header.html partial**

```html
{{define "partials/header"}}
<div class="flex items-center justify-between px-8 py-4">
    <div>
        <h2 class="text-2xl font-bold text-gray-800">{{.Title}}</h2>
        <p class="text-sm text-gray-500">{{.Subtitle}}</p>
    </div>

    <div class="flex items-center gap-4">
        <!-- Search bar -->
        <div class="relative">
            <input 
                type="text" 
                placeholder="Search..." 
                class="w-64 px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                hx-post="/search"
                hx-trigger="keyup changed delay:300ms"
                hx-target="#search-results"
            />
            <div id="search-results" class="absolute right-0 mt-2 w-96 bg-white rounded-lg shadow-lg z-20"></div>
        </div>

        <!-- User menu -->
        <div class="flex items-center gap-2">
            <span class="text-sm text-gray-600">{{.User.Name}}</span>
            <svg class="w-8 h-8 text-gray-400" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
        </div>
    </div>
</div>
{{end}}
```

- [ ] **Step 4: Create empty content template for testing**

Create `clinic-hcc-app/templates/pages/dashboard.html`:

```html
{{define "content"}}
<div class="card">
    <h3 class="text-lg font-semibold text-gray-800 mb-4">Welcome to Clinic HCC</h3>
    <p class="text-gray-600">Dashboard content coming soon.</p>
</div>
{{end}}

{{template "layouts/dashboard" .}}
```

- [ ] **Step 5: Commit**

```bash
git add clinic-hcc-app/templates/
git commit -m "feat: create base dashboard layout with sidebar and header"
```

---

### Task 5: Wire Up Templates and Test Server

**Files:**
- Modify: `clinic-hcc-app/internal/handlers/router.go`
- Modify: `clinic-hcc-app/internal/handlers/dashboard.go` (new file)

**Interfaces:**
- Consumes: Templates from Task 4, router from Task 2
- Produces: Working dashboard page

- [ ] **Step 1: Create dashboard.go handler**

```go
package handlers

import (
	"html/template"
	"net/http"
)

func (r *Router) Dashboard(w http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseGlob("templates/**/*.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Dashboard",
		"Subtitle":   "Overview",
		"ActivePage": "dashboard",
		"CSRFToken":  "TODO-generate-real-token",
		"User": map[string]string{
			"Name": "Practitioner",
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
```

- [ ] **Step 2: Update router.go to remove stub Dashboard method**

Remove the stub Dashboard method from router.go and keep only the signature. The implementation now lives in dashboard.go.

- [ ] **Step 3: Start TailwindCSS watch mode**

```bash
cd clinic-hcc-app
npm run dev:css
```

Expected: CSS rebuilds on file changes

- [ ] **Step 4: Start Go server in separate terminal**

```bash
cd clinic-hcc-app
go run ./cmd/server
```

Expected output:
```
2026/07/20 Starting server on :8080 (database: data/clinic.db)
```

- [ ] **Step 5: Test in browser**

Open http://localhost:8080

Expected: Dashboard page renders with:
- Sidebar navigation (left, fixed)
- Header with search bar and user menu (top)
- Main content area with welcome card
- TailwindCSS styles applied
- HTMX loaded (check browser console)

- [ ] **Step 6: Commit**

```bash
git add clinic-hcc-app/internal/handlers/dashboard.go clinic-hcc-app/internal/handlers/router.go
git commit -m "feat: wire up dashboard handler with template rendering"
```

---

### Task 6: Add Development Convenience Scripts

**Files:**
- Modify: `clinic-hcc-app/package.json`
- Create: `clinic-hcc-app/Makefile` (optional)

**Interfaces:**
- Consumes: Existing npm scripts
- Produces: Unified dev workflow

- [ ] **Step 1: Add dev script to package.json**

Update package.json scripts:
```json
{
  "scripts": {
    "dev:css": "npx tailwindcss -i src/index.css -o static/css/styles.css --watch",
    "build:css": "npx tailwindcss -i src/index.css -o static/css/styles.css --minify",
    "dev": "npm run dev:css"
  }
}
```

- [ ] **Step 2: Create Makefile for unified dev commands**

```makefile
.PHONY: dev build test clean

dev:
	@echo "Starting development servers..."
	@echo "Terminal 1: Go server"
	@echo "Terminal 2: TailwindCSS watch"
	@echo ""
	@echo "Run manually:"
	@echo "  Terminal 1: go run ./cmd/server"
	@echo "  Terminal 2: npm run dev:css"

build:
	go build -o server.exe ./cmd/server
	npm run build:css

test:
	go test ./...

clean:
	rm -f server.exe
	rm -f static/css/styles.css
```

- [ ] **Step 3: Update README with dev instructions**

Create `clinic-hcc-app/README.md`:

```markdown
# Clinic HCC App

Clinic receipt management system built with Go, HTMX, and TailwindCSS.

## Development

### Prerequisites

- Go 1.26.5+
- Node.js 18+
- SQLite 3

### Setup

```bash
# Install Go dependencies
go mod tidy

# Install Node dependencies
npm install
```

### Run Development Server

Two terminals required:

**Terminal 1 - Go server:**
```bash
go run ./cmd/server
```

**Terminal 2 - TailwindCSS watch:**
```bash
npm run dev:css
```

Open http://localhost:8080

### Build Production Binary

```bash
make build
# or
go build -o server.exe ./cmd/server
npm run build:css
```

### Run Tests

```bash
make test
# or
go test ./...
```
```

- [ ] **Step 4: Commit**

```bash
git add clinic-hcc-app/package.json clinic-hcc-app/Makefile clinic-hcc-app/README.md
git commit -m "docs: add development workflow documentation"
```

---

## Self-Review Checklist

**1. Spec Coverage:**
- [x] Chi router added and configured
- [x] HTTP routes match QUICKSTART.md specification
- [x] TailwindCSS v4 with CSS-first configuration (no tailwind.config.js)
- [x] Dashboard layout with sidebar and header
- [x] HTMX loaded from CDN
- [x] Development workflow documented

**2. Placeholder Scan:**
- [x] No TBD/TODO in task steps (except intentional stub implementations)
- [x] All code blocks contain actual implementation
- [x] All commands have expected outputs

**3. Type Consistency:**
- [x] Handler signatures consistent across router.go and individual handler files
- [x] Template data structures use map[string]interface{} consistently
- [x] File paths match actual project structure

---

## Execution Handoff

Plan complete and saved to `docs/superpowers/plans/2026-07-20-http-server-tailwind-setup.md`. Two execution options:

**1. Subagent-Driven (recommended)** - I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** - Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?