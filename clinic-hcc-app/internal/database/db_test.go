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

	if err := Migrate(db); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	tables := []string{"patients", "receipts", "receipt_items", "settings", "backups", "notifications"}
	for _, table := range tables {
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count); err != nil {
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

	if err := Migrate(db); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}

	var clinicName string
	if err := db.QueryRow("SELECT clinic_name FROM settings WHERE id = 1").Scan(&clinicName); err != nil {
		t.Fatalf("Failed to retrieve default settings: %v", err)
	}
	if clinicName != "Hong Ching Clinic" {
		t.Errorf("Expected clinic name 'Hong Ching Clinic', got %s", clinicName)
	}
}

func TestMigrateRepairsLegacyReceiptSchema(t *testing.T) {
	tmpFile := "test_legacy_migrate.db"
	defer os.Remove(tmpFile)

	db, err := New(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	legacySchema := `
		CREATE TABLE patients (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			hkid TEXT UNIQUE NOT NULL,
			hkid_hash TEXT NOT NULL,
			gender TEXT
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
			status TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (patient_id) REFERENCES patients(id)
		);
		CREATE TABLE receipt_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			receipt_id TEXT NOT NULL,
			description TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price INTEGER NOT NULL,
			subtotal INTEGER NOT NULL,
			FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE
		);
		INSERT INTO patients (id, name, hkid, hkid_hash, gender) VALUES ('p1', 'Patient', 'A123456(8)', 'hash', 'O');
		INSERT INTO receipts (id, receipt_number, patient_id, visit_date, subtotal, discount_type, grand_total, status)
		VALUES ('r1', 'RCP-OLD', 'p1', '2026-07-23', 100, 'none', 100, 'finalized');
		INSERT INTO receipt_items (receipt_id, description, quantity, unit_price, subtotal)
		VALUES ('r1', 'Treatment', 1, 100, 100);
	`
	if _, err := db.Exec(legacySchema); err != nil {
		t.Fatalf("create legacy schema: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate legacy schema: %v", err)
	}

	var notNull int
	if err := db.QueryRow(`SELECT "notnull" FROM pragma_table_info('receipts') WHERE name = 'receipt_number'`).Scan(&notNull); err != nil {
		t.Fatalf("inspect receipt number: %v", err)
	}
	if notNull != 0 {
		t.Fatalf("expected nullable receipt number, got notnull=%d", notNull)
	}
	var position int
	if err := db.QueryRow(`SELECT position FROM receipt_items WHERE receipt_id = 'r1'`).Scan(&position); err != nil {
		t.Fatalf("inspect migrated line item: %v", err)
	}
	if position != 0 {
		t.Fatalf("expected first migrated line item position 0, got %d", position)
	}
}
