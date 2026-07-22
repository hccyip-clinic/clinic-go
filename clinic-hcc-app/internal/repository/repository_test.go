package repository

import (
	"context"
	"path/filepath"
	"testing"

	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/models"
)

func testDB(t *testing.T) *database.DB {
	t.Helper()
	db, err := database.New(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("create database: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		db.Close()
		t.Fatalf("migrate database: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestPatientRepository_NormalizesAndRejectsDuplicateHKID(t *testing.T) {
	db := testDB(t)
	repo := NewPatientRepository(db)
	ctx := context.Background()

	patient := &models.Patient{Name: " Chan Tai Man ", HKID: "a 123-456(8)", Gender: "O"}
	if err := repo.Create(ctx, patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	if patient.HKID != "A123456(8)" {
		t.Fatalf("expected canonical HKID, got %q", patient.HKID)
	}

	duplicate := &models.Patient{Name: "Duplicate", HKID: "A123456(8)", Gender: "M"}
	if err := repo.Create(ctx, duplicate); err == nil {
		t.Fatal("expected duplicate HKID to be rejected")
	}
}

func TestReceiptRepository_FinalizeMakesReceiptImmutable(t *testing.T) {
	db := testDB(t)
	patients := NewPatientRepository(db)
	receipts := NewReceiptRepository(db)
	ctx := context.Background()

	patient := &models.Patient{Name: "Chan Tai Man", HKID: "A123456(8)", Gender: "O"}
	if err := patients.Create(ctx, patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	receipt := &models.Receipt{
		PatientID:    patient.ID,
		VisitDate:    "2026-07-22",
		DiscountType: models.DiscountNone,
		LineItems: []models.LineItem{{
			Description: "Tui Na",
			Quantity:    1,
			UnitPrice:   60000,
		}},
	}
	if err := receipts.CreateDraft(ctx, receipt); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	if receipt.ReceiptNumber != "" {
		t.Fatal("draft should not have a receipt number")
	}
	if err := receipts.Finalize(ctx, receipt.ID, "RCP"); err != nil {
		t.Fatalf("finalize receipt: %v", err)
	}
	finalized, err := receipts.Get(ctx, receipt.ID)
	if err != nil {
		t.Fatalf("get finalized receipt: %v", err)
	}
	if finalized.Status != models.StatusFinalized || finalized.ReceiptNumber == "" {
		t.Fatalf("expected finalized receipt with number, got %#v", finalized)
	}
	if err := receipts.DeleteDraft(ctx, receipt.ID); err != ErrReceiptImmutable {
		t.Fatalf("expected immutable error, got %v", err)
	}
}

func TestReceiptRepository_PreservesLineItemOrder(t *testing.T) {
	db := testDB(t)
	patients := NewPatientRepository(db)
	receipts := NewReceiptRepository(db)
	ctx := context.Background()

	patient := &models.Patient{Name: "Chan Tai Man", HKID: "A123456(8)", Gender: "O"}
	if err := patients.Create(ctx, patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}

	receipt := &models.Receipt{
		PatientID:    patient.ID,
		VisitDate:    "2026-07-23",
		DiscountType: models.DiscountNone,
		LineItems: []models.LineItem{
			{Description: "Second entered", Quantity: 1, UnitPrice: 200},
			{Description: "First entered", Quantity: 1, UnitPrice: 100},
		},
	}

	if err := receipts.CreateDraft(ctx, receipt); err != nil {
		t.Fatalf("create draft: %v", err)
	}
	loaded, err := receipts.Get(ctx, receipt.ID)
	if err != nil {
		t.Fatalf("get draft: %v", err)
	}
	if len(loaded.LineItems) != 2 || loaded.LineItems[0].Description != "Second entered" || loaded.LineItems[1].Description != "First entered" {
		t.Fatalf("line-item order was not preserved: %#v", loaded.LineItems)
	}
}

func TestReceiptRepository_RejectsInvalidDraft(t *testing.T) {
	db := testDB(t)
	patients := NewPatientRepository(db)
	receipts := NewReceiptRepository(db)
	ctx := context.Background()

	patient := &models.Patient{Name: "Chan Tai Man", HKID: "A123456(8)", Gender: "O"}
	if err := patients.Create(ctx, patient); err != nil {
		t.Fatalf("create patient: %v", err)
	}
	err := receipts.CreateDraft(ctx, &models.Receipt{
		PatientID:    patient.ID,
		VisitDate:    "2026-07-23",
		DiscountType: models.DiscountNone,
	})
	if err == nil {
		t.Fatal("expected empty draft to be rejected")
	}
}
