package database

import (
	"database/sql"
	"fmt"
	"strings"
)

const schemaSQL = `
-- Patients table
CREATE TABLE IF NOT EXISTS patients (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	hkid TEXT UNIQUE NOT NULL,
	hkid_hash TEXT NOT NULL,
	gender TEXT CHECK(gender IN ('M', 'F', 'O')),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Receipts table
CREATE TABLE IF NOT EXISTS receipts (
	id TEXT PRIMARY KEY,
	receipt_number TEXT UNIQUE,
	patient_id TEXT NOT NULL,
	visit_date DATE NOT NULL,
	diagnosis TEXT,
	subtotal INTEGER NOT NULL,
	discount_type TEXT CHECK(discount_type IN ('percent', 'fixed', 'none')),
	discount_value INTEGER DEFAULT 0,
	grand_total INTEGER NOT NULL,
	status TEXT CHECK(status IN ('draft', 'finalized', 'archived')),
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (patient_id) REFERENCES patients(id)
);

-- Receipt line items
CREATE TABLE IF NOT EXISTS receipt_items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	receipt_id TEXT NOT NULL,
	position INTEGER NOT NULL,
	description TEXT NOT NULL,
	quantity INTEGER NOT NULL,
	unit_price INTEGER NOT NULL,
	subtotal INTEGER NOT NULL,
	FOREIGN KEY (receipt_id) REFERENCES receipts(id) ON DELETE CASCADE
);

-- Clinic settings (single row)
CREATE TABLE IF NOT EXISTS settings (
	id INTEGER PRIMARY KEY CHECK(id = 1),
	clinic_name TEXT NOT NULL,
	clinic_address TEXT,
	clinic_phone TEXT,
	practitioner_name TEXT NOT NULL,
	practitioner_registration TEXT,
	receipt_prefix TEXT DEFAULT 'RCP',
	retention_years INTEGER DEFAULT 3,
	password_hash TEXT NOT NULL DEFAULT '',
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Auto-backups log
CREATE TABLE IF NOT EXISTS backups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	backup_type TEXT CHECK(backup_type IN ('delta', 'full')),
	backup_date DATE NOT NULL,
	file_path TEXT,
	masked INTEGER DEFAULT 1,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Notifications table (7-day retention)
CREATE TABLE IF NOT EXISTS notifications (
	id TEXT PRIMARY KEY,
	category TEXT NOT NULL,
	scope TEXT NOT NULL,
	title TEXT NOT NULL,
	message TEXT NOT NULL,
	action_url TEXT,
	is_read INTEGER DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	expires_at DATETIME
);

-- Weekly metrics cache (optional optimization)
CREATE TABLE IF NOT EXISTS weekly_metrics_cache (
	date DATE PRIMARY KEY,
	patient_count INTEGER,
	receipt_count INTEGER,
	total_amount INTEGER
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_receipts_patient ON receipts(patient_id);
CREATE INDEX IF NOT EXISTS idx_receipts_visit_date ON receipts(visit_date);
CREATE INDEX IF NOT EXISTS idx_receipts_status ON receipts(status);
CREATE INDEX IF NOT EXISTS idx_patients_hkid ON patients(hkid);
CREATE INDEX IF NOT EXISTS idx_notifications_expires ON notifications(expires_at);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);

CREATE TABLE IF NOT EXISTS schema_migrations (
	version TEXT PRIMARY KEY,
	applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

func Migrate(db *DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return err
	}
	if err := migrateReceiptSchema(db); err != nil {
		return err
	}
	// Keep existing Phase 1 databases upgradeable.
	if _, err := db.Exec(`ALTER TABLE settings ADD COLUMN password_hash TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(err.Error(), "duplicate column") {
		return err
	}

	return SeedDefaults(db)
}

func migrateReceiptSchema(db *DB) error {
	const version = "receipt-schema-v2"
	var applied int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = ?`, version).Scan(&applied); err != nil {
		return err
	}
	if applied > 0 {
		return nil
	}

	receiptNumberNotNull, err := columnNotNull(db, "receipts", "receipt_number")
	if err != nil {
		return err
	}
	positionExists, err := columnExists(db, "receipt_items", "position")
	if err != nil {
		return err
	}
	if receiptNumberNotNull || !positionExists {
		if err := rebuildReceiptTables(db, receiptNumberNotNull, positionExists); err != nil {
			return err
		}
	}
	_, err = db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version)
	return err
}

func columnExists(db *DB, table, column string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return true, nil
		}
	}
	return false, rows.Err()
}

func columnNotNull(db *DB, table, column string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, dataType string
		var notNull, primaryKey int
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if name == column {
			return notNull == 1, nil
		}
	}
	return false, rows.Err()
}

func rebuildReceiptTables(db *DB, _, positionExists bool) (err error) {
	if _, err := db.Exec(`PRAGMA foreign_keys = OFF`); err != nil {
		return err
	}
	defer func() {
		if restoreErr := func() error {
			_, restoreErr := db.Exec(`PRAGMA foreign_keys = ON`)
			return restoreErr
		}(); restoreErr != nil && err == nil {
			err = fmt.Errorf("restore foreign keys after migration: %w", restoreErr)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var receiptCount, itemCount int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM receipts`).Scan(&receiptCount); err != nil {
		return err
	}
	if err := tx.QueryRow(`SELECT COUNT(*) FROM receipt_items`).Scan(&itemCount); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		CREATE TABLE receipts_new (
			id TEXT PRIMARY KEY,
			receipt_number TEXT UNIQUE,
			patient_id TEXT NOT NULL,
			visit_date DATE NOT NULL,
			diagnosis TEXT,
			subtotal INTEGER NOT NULL,
			discount_type TEXT CHECK(discount_type IN ('percent', 'fixed', 'none')),
			discount_value INTEGER DEFAULT 0,
			grand_total INTEGER NOT NULL,
			status TEXT CHECK(status IN ('draft', 'finalized', 'archived')),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (patient_id) REFERENCES patients(id)
		)
	`); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		INSERT INTO receipts_new (
			id, receipt_number, patient_id, visit_date, diagnosis, subtotal,
			discount_type, discount_value, grand_total, status, created_at, updated_at
		)
		SELECT id, receipt_number, patient_id, visit_date, diagnosis, subtotal,
			discount_type, discount_value, grand_total, status, created_at, updated_at
		FROM receipts
	`); err != nil {
		return err
	}
	if _, err := tx.Exec(`
		CREATE TABLE receipt_items_new (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			receipt_id TEXT NOT NULL,
			position INTEGER NOT NULL,
			description TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price INTEGER NOT NULL,
			subtotal INTEGER NOT NULL,
			FOREIGN KEY (receipt_id) REFERENCES receipts_new(id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}
	positionExpression := "ROW_NUMBER() OVER (PARTITION BY receipt_id ORDER BY position, id) - 1"
	if !positionExists {
		positionExpression = "ROW_NUMBER() OVER (PARTITION BY receipt_id ORDER BY id) - 1"
	}
	if _, err := tx.Exec(`
		INSERT INTO receipt_items_new (id, receipt_id, position, description, quantity, unit_price, subtotal)
		SELECT id, receipt_id, ` + positionExpression + `, description, quantity, unit_price, subtotal
		FROM receipt_items
	`); err != nil {
		return err
	}
	var migratedReceiptCount, migratedItemCount int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM receipts_new`).Scan(&migratedReceiptCount); err != nil {
		return err
	}
	if err := tx.QueryRow(`SELECT COUNT(*) FROM receipt_items_new`).Scan(&migratedItemCount); err != nil {
		return err
	}
	if migratedReceiptCount != receiptCount || migratedItemCount != itemCount {
		return fmt.Errorf("receipt migration row count mismatch: receipts %d/%d, items %d/%d", migratedReceiptCount, receiptCount, migratedItemCount, itemCount)
	}
	if _, err := tx.Exec(`DROP TABLE receipt_items; DROP TABLE receipts;`); err != nil {
		return err
	}
	if _, err := tx.Exec(`ALTER TABLE receipts_new RENAME TO receipts; ALTER TABLE receipt_items_new RENAME TO receipt_items;`); err != nil {
		return err
	}
	foreignKeyRows, err := tx.Query(`PRAGMA foreign_key_check`)
	if err != nil {
		return err
	}
	defer foreignKeyRows.Close()
	if foreignKeyRows.Next() {
		return fmt.Errorf("receipt migration produced a foreign-key violation")
	}
	if err := foreignKeyRows.Err(); err != nil {
		return err
	}
	for _, statement := range []string{
		`CREATE INDEX IF NOT EXISTS idx_receipts_patient ON receipts(patient_id)`,
		`CREATE INDEX IF NOT EXISTS idx_receipts_visit_date ON receipts(visit_date)`,
		`CREATE INDEX IF NOT EXISTS idx_receipts_status ON receipts(status)`,
	} {
		if _, err := tx.Exec(statement); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func SeedDefaults(db *DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM settings").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = db.Exec(`
			INSERT INTO settings (
				id, clinic_name, clinic_address, clinic_phone,
				practitioner_name, practitioner_registration, receipt_prefix, retention_years
			) VALUES (
				1, 'Hong Ching Clinic', '', '',
				'Practitioner', '', 'RCP', 3
			)
		`)
		if err != nil {
			return err
		}
	}

	return nil
}
