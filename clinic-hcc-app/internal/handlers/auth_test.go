package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"clinic-hcc-app/internal/database"
)

func TestUnauthenticatedRequestsRedirectToSetup(t *testing.T) {
	path := "auth_handler_test.db"
	defer os.Remove(path)
	db, err := database.New(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	NewRouter(db).Setup().ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/setup" {
		t.Fatalf("expected setup redirect, got %d %q", rec.Code, rec.Header().Get("Location"))
	}
}
