package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/models"
)

var (
	ErrReceiptNotFound  = errors.New("receipt not found")
	ErrReceiptImmutable = errors.New("finalized receipts are immutable")
)

type ReceiptRepository struct {
	db *database.DB
}

func NewReceiptRepository(db *database.DB) *ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) CreateDraft(ctx context.Context, receipt *models.Receipt) error {
	if receipt.ID == "" {
		var err error
		receipt.ID, err = newID("receipt")
		if err != nil {
			return err
		}
	}
	receipt.Status = models.StatusDraft
	receipt.Subtotal = models.CalculateSubtotal(receipt.LineItems)
	receipt.GrandTotal = models.CalculateGrandTotal(receipt.Subtotal, receipt.DiscountType, receipt.DiscountValue)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO receipts (id, receipt_number, patient_id, visit_date, diagnosis,
			subtotal, discount_type, discount_value, grand_total, status)
		VALUES (?, NULL, ?, ?, ?, ?, ?, ?, ?, ?)
	`, receipt.ID, receipt.PatientID, receipt.VisitDate, receipt.Diagnosis,
		receipt.Subtotal, receipt.DiscountType, receipt.DiscountValue, receipt.GrandTotal, receipt.Status)
	if err != nil {
		return err
	}
	if err := insertItems(ctx, tx, receipt.ID, receipt.LineItems); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ReceiptRepository) Get(ctx context.Context, id string) (*models.Receipt, error) {
	receipt := &models.Receipt{}
	err := r.db.QueryRowContext(ctx, `
		SELECT id, COALESCE(receipt_number, ''), patient_id, visit_date, diagnosis,
			subtotal, discount_type, discount_value, grand_total, status, created_at, updated_at
		FROM receipts WHERE id = ?
	`, id).Scan(
		&receipt.ID, &receipt.ReceiptNumber, &receipt.PatientID, &receipt.VisitDate,
		&receipt.Diagnosis, &receipt.Subtotal, &receipt.DiscountType, &receipt.DiscountValue,
		&receipt.GrandTotal, &receipt.Status, &receipt.CreatedAt, &receipt.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrReceiptNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, description, quantity, unit_price, subtotal
		FROM receipt_items WHERE receipt_id = ? ORDER BY id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.LineItem
		if err := rows.Scan(&item.ID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Subtotal); err != nil {
			return nil, err
		}
		receipt.LineItems = append(receipt.LineItems, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return receipt, nil
}

func (r *ReceiptRepository) List(ctx context.Context) ([]models.Receipt, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, COALESCE(receipt_number, ''), patient_id, visit_date, diagnosis,
			subtotal, discount_type, discount_value, grand_total, status, created_at, updated_at
		FROM receipts ORDER BY visit_date DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var receipts []models.Receipt
	for rows.Next() {
		var receipt models.Receipt
		if err := rows.Scan(
			&receipt.ID, &receipt.ReceiptNumber, &receipt.PatientID, &receipt.VisitDate,
			&receipt.Diagnosis, &receipt.Subtotal, &receipt.DiscountType, &receipt.DiscountValue,
			&receipt.GrandTotal, &receipt.Status, &receipt.CreatedAt, &receipt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, rows.Err()
}

func (r *ReceiptRepository) UpdateDraft(ctx context.Context, receipt *models.Receipt) error {
	current, err := r.Get(ctx, receipt.ID)
	if err != nil {
		return err
	}
	if current.Status != models.StatusDraft {
		return ErrReceiptImmutable
	}
	receipt.Subtotal = models.CalculateSubtotal(receipt.LineItems)
	receipt.GrandTotal = models.CalculateGrandTotal(receipt.Subtotal, receipt.DiscountType, receipt.DiscountValue)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `
		UPDATE receipts SET patient_id = ?, visit_date = ?, diagnosis = ?,
			subtotal = ?, discount_type = ?, discount_value = ?, grand_total = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'draft'
	`, receipt.PatientID, receipt.VisitDate, receipt.Diagnosis, receipt.Subtotal,
		receipt.DiscountType, receipt.DiscountValue, receipt.GrandTotal, receipt.ID)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM receipt_items WHERE receipt_id = ?`, receipt.ID); err != nil {
		return err
	}
	if err := insertItems(ctx, tx, receipt.ID, receipt.LineItems); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ReceiptRepository) DeleteDraft(ctx context.Context, id string) error {
	receipt, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if receipt.Status != models.StatusDraft {
		return ErrReceiptImmutable
	}
	_, err = r.db.ExecContext(ctx, `DELETE FROM receipts WHERE id = ? AND status = 'draft'`, id)
	return err
}

func (r *ReceiptRepository) Finalize(ctx context.Context, id, prefix string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	receipt, err := getReceiptTx(ctx, tx, id)
	if err != nil {
		return err
	}
	if receipt.Status != models.StatusDraft {
		return ErrReceiptImmutable
	}
	items, err := getItemsTx(ctx, tx, id)
	if err != nil {
		return err
	}
	receipt.LineItems = items
	receipt.Subtotal = models.CalculateSubtotal(items)
	receipt.GrandTotal = models.CalculateGrandTotal(receipt.Subtotal, receipt.DiscountType, receipt.DiscountValue)
	if validationErrors := models.ValidateReceipt(receipt); len(validationErrors) > 0 {
		return fmt.Errorf("receipt validation failed: %v", validationErrors)
	}

	var number string
	for attempts := 0; attempts < 10; attempts++ {
		number, err = models.GenerateReceiptNumber(prefix)
		if err != nil {
			return err
		}
		var count int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM receipts WHERE receipt_number = ?`, number).Scan(&count); err != nil {
			return err
		}
		if count == 0 {
			break
		}
		number = ""
	}
	if number == "" {
		return errors.New("could not generate unique receipt number")
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE receipts SET receipt_number = ?, subtotal = ?, grand_total = ?,
			status = 'finalized', updated_at = CURRENT_TIMESTAMP
		WHERE id = ? AND status = 'draft'
	`, number, receipt.Subtotal, receipt.GrandTotal, id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func insertItems(ctx context.Context, tx *sql.Tx, receiptID string, items []models.LineItem) error {
	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO receipt_items (receipt_id, description, quantity, unit_price, subtotal)
			VALUES (?, ?, ?, ?, ?)
		`, receiptID, item.Description, item.Quantity, item.UnitPrice, item.Quantity*item.UnitPrice); err != nil {
			return err
		}
	}
	return nil
}

func getReceiptTx(ctx context.Context, tx *sql.Tx, id string) (*models.Receipt, error) {
	receipt := &models.Receipt{}
	err := tx.QueryRowContext(ctx, `
		SELECT id, COALESCE(receipt_number, ''), patient_id, visit_date, diagnosis,
			subtotal, discount_type, discount_value, grand_total, status, created_at, updated_at
		FROM receipts WHERE id = ?
	`, id).Scan(
		&receipt.ID, &receipt.ReceiptNumber, &receipt.PatientID, &receipt.VisitDate,
		&receipt.Diagnosis, &receipt.Subtotal, &receipt.DiscountType, &receipt.DiscountValue,
		&receipt.GrandTotal, &receipt.Status, &receipt.CreatedAt, &receipt.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrReceiptNotFound
	}
	return receipt, err
}

func getItemsTx(ctx context.Context, tx *sql.Tx, receiptID string) ([]models.LineItem, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, description, quantity, unit_price, subtotal
		FROM receipt_items WHERE receipt_id = ? ORDER BY id
	`, receiptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []models.LineItem
	for rows.Next() {
		var item models.LineItem
		if err := rows.Scan(&item.ID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Subtotal); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
