package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/models"
)

var ErrPatientNotFound = errors.New("patient not found")

type PatientRepository struct {
	db *database.DB
}

func NewPatientRepository(db *database.DB) *PatientRepository {
	return &PatientRepository{db: db}
}

func (r *PatientRepository) Create(ctx context.Context, patient *models.Patient) error {
	canonical, err := models.NormalizeHKID(patient.HKID)
	if err != nil {
		return err
	}
	patient.HKID = canonical
	patient.Name = strings.TrimSpace(patient.Name)
	if patient.Name == "" {
		return errors.New("patient name is required")
	}
	if patient.Gender != "M" && patient.Gender != "F" && patient.Gender != "O" {
		return errors.New("patient gender must be M, F, or O")
	}
	patient.HKIDHash = models.HashHKID(canonical)
	if patient.ID == "" {
		patient.ID, err = newID("patient")
		if err != nil {
			return err
		}
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO patients (id, name, hkid, hkid_hash, gender)
		VALUES (?, ?, ?, ?, ?)
	`, patient.ID, patient.Name, patient.HKID, patient.HKIDHash, patient.Gender)
	return err
}

func (r *PatientRepository) Get(ctx context.Context, id string) (*models.Patient, error) {
	patient := &models.Patient{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, hkid, hkid_hash, gender, created_at, updated_at
		FROM patients WHERE id = ?
	`, id).Scan(
		&patient.ID, &patient.Name, &patient.HKID, &patient.HKIDHash,
		&patient.Gender, &patient.CreatedAt, &patient.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPatientNotFound
	}

	if err != nil {
		return nil, err
	}
	return patient, nil
}

func (r *PatientRepository) Update(ctx context.Context, patient *models.Patient) error {
	patient.Name = strings.TrimSpace(patient.Name)
	if patient.Name == "" {
		return errors.New("patient name is required")
	}
	if patient.Gender != "M" && patient.Gender != "F" && patient.Gender != "O" {
		return errors.New("patient gender must be M, F, or O")
	}
	result, err := r.db.ExecContext(ctx, `
		UPDATE patients SET name = ?, gender = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, patient.Name, patient.Gender, patient.ID)
	if err != nil {
		return err
	}
	if count, err := result.RowsAffected(); err != nil {
		return err
	} else if count == 0 {
		return ErrPatientNotFound
	}
	return nil
}

func (r *PatientRepository) Search(ctx context.Context, query string) ([]models.Patient, error) {
	query = strings.TrimSpace(query)
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, hkid, hkid_hash, gender, created_at, updated_at
		FROM patients
		WHERE name LIKE '%' || ? || '%' COLLATE NOCASE
		   OR hkid = ?
		ORDER BY name, created_at DESC
	`, query, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []models.Patient
	for rows.Next() {
		var patient models.Patient
		if err := rows.Scan(
			&patient.ID, &patient.Name, &patient.HKID, &patient.HKIDHash,
			&patient.Gender, &patient.CreatedAt, &patient.UpdatedAt,
		); err != nil {
			return nil, err
		}
		patients = append(patients, patient)
	}
	return patients, rows.Err()
}

func newID(prefix string) (string, error) {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate %s id: %w", prefix, err)
	}
	return fmt.Sprintf("%s-%x", prefix, buffer), nil
}
