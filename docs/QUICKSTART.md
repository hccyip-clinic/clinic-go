# Quick-Start Guide: Hugo + HTMX + SQLite Clinic App

## Prerequisites

- Go 1.21 or later
- SQLite 3
- Node.js 18+ (for TailwindCSS CLI)
- `air` (optional, for hot reload)

## Project Setup

### 1. Initialize Go Module

```bash
mkdir clinic-app
cd clinic-app
go mod init clinic-app
```

### 2. Install Dependencies

```bash
# Go dependencies
go get github.com/mattn/go-sqlite3
go get github.com/go-chi/chi/v5

# Optional: hot reload for development
go install github.com/cosmtrek/air@latest

# TailwindCSS
npm install -D tailwindcss @tailwindcss/vite
```

### 3. Directory Structure

```
clinic-app/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── handlers/             # HTTP handlers
│   ├── models/               # Go structs
│   ├── repository/           # SQLite queries
│   └── services/             # Business logic
├── templates/                # HTML templates
│   ├── layouts/
│   │   └── dashboard.html    # Base layout
│   ├── pages/
│   │   ├── dashboard.html
│   │   ├── receipts.html
│   │   ├── receipt-form.html
│   │   ├── patients.html
│   │   └── settings.html
│   └── partials/
│       ├── sidebar.html
│       ├── header.html
│       └── receipt-row.html
├── static/
│   ├── css/
│   │   └── styles.css        # Tailwind output
│   └── js/
│       └── htmx.min.js       # HTMX CDN or local
├── data/
│   └── clinic.db             # SQLite database (gitignore)
├── input.css                 # Tailwind input
├── tailwind.config.js
└── go.mod
```

---

## Database Setup

### Create Schema

```bash
sqlite3 data/clinic.db < schema.sql
```

**schema.sql:**
```sql
-- See docs/prototype-spec.md for full schema
CREATE TABLE patients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    hkid TEXT UNIQUE NOT NULL,
    hkid_hash TEXT NOT NULL,
    gender TEXT CHECK(gender IN ('M', 'F', 'O')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receipts (
    id TEXT PRIMARY KEY,
    receipt_number TEXT UNIQUE NOT NULL,
    patient_id TEXT NOT NULL,
    visit_date DATE NOT NULL,
    diagnosis TEXT,
    subtotal INTEGER NOT NULL,
    discount_type TEXT,
    discount_value INTEGER DEFAULT 0,
    grand_total INTEGER NOT NULL,
    status TEXT CHECK(status IN ('draft', 'finalized', 'archived')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (patient_id) REFERENCES patients(id)
);

-- Add remaining tables from docs/prototype-spec.md
```

---

## HTTP Server

### main.go

```go
package main

import (
    "clinic-app/internal/handlers"
    "clinic-app/internal/repository"
    "log"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    // Initialize database
    db, err := repository.NewSQLite("data/clinic.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Create router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

    // Serve static files
    r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // Routes
    h := handlers.NewHandler(db)
    r.Get("/", h.Dashboard)
    r.Get("/receipts", h.ReceiptList)
    r.Get("/receipts/new", h.ReceiptFormNew)
    r.Post("/receipts", h.ReceiptCreate)
    r.Get("/receipts/{id}", h.ReceiptView)
    r.Get("/receipts/{id}/edit", h.ReceiptFormEdit)
    r.Post("/receipts/{id}", h.ReceiptUpdate)
    r.Delete("/receipts/{id}", h.ReceiptDelete)
    r.Get("/patients", h.PatientList)
    r.Get("/settings", h.Settings)
    r.Post("/settings", h.SettingsUpdate)

    // Start server
    log.Println("Starting server on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

---

## Handlers Example

### handlers/receipt.go

```go
package handlers

import (
    "clinic-app/internal/repository"
    "database/sql"
    "html/template"
    "net/http"
    "strconv"
)

type ReceiptHandler struct {
    db    *sql.DB
    tmpl  *template.Template
}

func NewHandler(db *sql.DB) *ReceiptHandler {
    return &ReceiptHandler{
        db:   db,
        tmpl: template.Must(template.ParseGlob("templates/**/*.html")),
    }
}

func (h *ReceiptHandler) ReceiptList(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page == 0 {
        page = 1
    }

    receipts, err := repository.ListReceipts(h.db, page, 20)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    data := map[string]interface{}{
        "Receipts": receipts,
        "Page":     page,
    }

    h.tmpl.ExecuteTemplate(w, "pages/receipts.html", data)
}

func (h *ReceiptHandler) ReceiptCreate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", 405)
        return
    }

    // Parse form data
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    // Validate and create receipt
    // ... validation logic here ...

    // For HTMX partial update, return just the form fragment with errors
    // For full page, redirect
    if r.Header.Get("HX-Request") == "true" {
        // Return validation errors as fragment
        h.tmpl.ExecuteTemplate(w, "partials/receipt-form-errors.html", data)
    } else {
        http.Redirect(w, r, "/receipts", http.StatusSeeOther)
    }
}
```

---

## HTMX Form Example

### templates/pages/receipt-form.html

```html
{{define "page-title"}}New Receipt{{end}}

{{define "content"}}
<form 
    hx-post="/receipts"
    hx-swap="outerHTML"
    hx-target="#receipt-form"
    hx-on::after-request="if(event.detail.successful) window.location='/receipts'"
    id="receipt-form"
>
    <!-- Patient Info -->
    <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700">Patient Name</label>
        <input 
            type="text" 
            name="patientName" 
            required
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            hx-post="/patients/search"
            hx-trigger="keyup changed delay:300ms"
            hx-target="#patient-suggestions"
        />
        <div id="patient-suggestions"></div>
    </div>

    <!-- HKID with validation -->
    <div class="mb-6">
        <label class="block text-sm font-medium text-gray-700">HKID</label>
        <input 
            type="text" 
            name="hkid" 
            required
            pattern="[A-Z]{1,2}[0-9]{6}\([0-9A]\)"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm"
        />
        <div class="text-sm text-red-600" id="hkid-error"></div>
    </div>

    <!-- Line Items (dynamic) -->
    <div id="line-items">
        {{range .LineItems}}
        <div class="line-item mb-4 p-4 border rounded">
            <input type="text" name="description" placeholder="Description" class="w-full" />
            <input type="number" name="quantity" value="1" class="w-20" />
            <input type="number" name="unitPrice" placeholder="Price" class="w-32" />
            <button type="button" 
                    hx-delete="/receipts/items/{id}"
                    hx-swap="outerHTML"
                    class="text-red-600 hover:text-red-800">
                Remove
            </button>
        </div>
        {{end}}
    </div>

    <button type="button"
            hx-post="/receipts/items/new"
            hx-swap="beforeend"
            hx-target="#line-items"
            class="mb-4 text-blue-600 hover:text-blue-800">
        + Add Line Item
    </button>

    <!-- Totals (auto-calculated) -->
    <div class="mt-6">
        <div>Subtotal: $<span id="subtotal">{{.Subtotal}}</span></div>
        <div>Discount: $<span id="discount">{{.Discount}}</span></div>
        <div class="text-xl font-bold">Total: $<span id="grand-total">{{.GrandTotal}}</span></div>
    </div>

    <!-- Actions -->
    <div class="mt-8 flex gap-4">
        <button type="submit" name="action" value="draft" class="px-4 py-2 bg-gray-200 rounded">
            Save Draft
        </button>
        <button type="submit" name="action" value="finalize" class="px-4 py-2 bg-blue-600 text-white rounded">
            Finalize
        </button>
    </div>
</form>
{{end}}
```

---

## TailwindCSS Setup

### input.css

```css
@import "tailwindcss";

@theme {
  --color-primary: #3b82f6;
  --color-secondary: #64748b;
  --font-sans: 'Inter', system-ui, sans-serif;
  --radius-lg: 0.75rem;
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
```

### tailwind.config.js

```javascript
/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./templates/**/*.html",
    "./static/js/**/*.js",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
```

### Build Command

```bash
npx tailwindcss -i input.css -o static/css/styles.css --watch
```

---

## Dashboard Layout

### templates/layouts/dashboard.html

```html
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
    <aside class="sidebar-fixed bg-white shadow-lg">
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
```

---

## Development Workflow

### Run Server

```bash
# Terminal 1: Go server
go run cmd/server/main.go

# Terminal 2: Tailwind watch
npx tailwindcss -i input.css -o static/css/styles.css --watch

# Optional: Hot reload with air
air
```

### Access Application

Open http://localhost:8080 in your browser.

---

## Testing

### Run Tests

```bash
go test ./...
```

### Manual Testing Checklist

- [ ] Create new patient
- [ ] Create new receipt with line items
- [ ] Edit existing receipt
- [ ] Delete receipt
- [ ] Search patients by HKID
- [ ] Generate monthly report
- [ ] Export CSV
- [ ] Settings update

---

## Deployment

### Build Binary

```bash
CGO_ENABLED=1 go build -o clinic-app cmd/server/main.go
```

### Run Production

```bash
./clinic-app
```

The application runs as a single binary with embedded SQLite. No external dependencies required.

---

## Next Steps

1. **Read full specification**: `docs/prototype-spec.md`
2. **Review domain model**: `CONTEXT.md`
3. **Understand trade-offs**: `docs/hugo-htmx-sqlite.md`
4. **Start prototyping**: Follow Phase 1 in prototype-spec.md

---

## Common Issues

### "SQLite database is locked"

Enable WAL mode in schema:
```sql
PRAGMA journal_mode=WAL;
```

### HTMX not triggering

Check:
- HTMX script loaded (`<script src="...htmx...">`)
- `hx-post` attribute spelled correctly
- Server returns appropriate status code (200, 204, 422)

### Tailwind classes not applying

- Verify `content` array in `tailwind.config.js` includes template paths
- Run `npx tailwindcss --watch` to rebuild CSS
- Clear browser cache

---

## Resources

- **HTMX Reference**: https://htmx.org/reference/
- **TailwindCSS v4**: https://tailwindcss.com/docs
- **Go SQLite**: https://github.com/mattn/go-sqlite3
- **Chi Router**: https://github.com/go-chi/chi