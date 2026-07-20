package database

import (
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tmpFile := "test.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if db.DB == nil {
		t.Error("Database connection is nil")
	}

	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode").Scan(&journalMode)
	if err != nil {
		t.Fatalf("Failed to check journal mode: %v", err)
	}
	if journalMode != "wal" {
		t.Errorf("Expected WAL mode, got %s", journalMode)
	}

	var foreignKeys int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&foreignKeys)
	if err != nil {
		t.Fatalf("Failed to check foreign keys: %v", err)
	}
	if foreignKeys != 1 {
		t.Error("Foreign keys not enabled")
	}
}

func TestMigrate(t *testing.T) {
	tmpFile := "test_migrate.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = Migrate(db)
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	tables := []string{"patients", "receipts", "receipt_items", "settings", "backups", "notifications"}
	for _, table := range tables {
		var count int
		query := "SELECT COUNT(*) FROM " + table
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			t.Errorf("Table %s does not exist: %v", table, err)
		}
	}
}

func TestSeedDefaults(t *testing.T) {
	tmpFile := "test_seed.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = Migrate(db)
	if err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	var clinicName string
	err = db.QueryRow("SELECT clinic_name FROM settings WHERE id = 1").Scan(&clinicName)
	if err != nil {
		t.Fatalf("Failed to retrieve default settings: %v", err)
	}

	if clinicName != "Hong Ching Clinic" {
		t.Errorf("Expected clinic name 'Hong Ching Clinic', got %s", clinicName)
	}
}