package database

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
`

func Migrate(db *DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return err
	}

	return SeedDefaults(db)
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
