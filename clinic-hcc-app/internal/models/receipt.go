package models

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"
)

type ReceiptStatus string

const (
	StatusDraft     ReceiptStatus = "draft"
	StatusFinalized ReceiptStatus = "finalized"
	StatusArchived  ReceiptStatus = "archived"
)

type DiscountType string

const (
	DiscountNone    DiscountType = "none"
	DiscountPercent DiscountType = "percent"
	DiscountFixed   DiscountType = "fixed"
)

type LineItem struct {
	ID          int64
	Description string
	Quantity    int
	UnitPrice   int
	Subtotal    int
}

type Receipt struct {
	ID            string
	ReceiptNumber string
	PatientID     string
	VisitDate     string
	Diagnosis     string
	LineItems     []LineItem
	Subtotal      int
	DiscountType  DiscountType
	DiscountValue int
	GrandTotal    int
	Status        ReceiptStatus
	CreatedAt     string
	UpdatedAt     string
}

func GenerateReceiptNumber(prefix string) (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, 6)
	buffer := make([]byte, len(random))
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("generate receipt number: %w", err)
	}
	for i := range random {
		random[i] = chars[int(buffer[i])%len(chars)]
	}

	return fmt.Sprintf("%s-%s-%s", strings.TrimSpace(prefix), time.Now().Format("20060102"), string(random)), nil
}

func CalculateSubtotal(items []LineItem) int {
	total := 0
	for _, item := range items {
		total += item.Quantity * item.UnitPrice
	}
	return total
}

func CalculateGrandTotal(subtotal int, discountType DiscountType, discountValue int) int {
	switch discountType {
	case DiscountNone:
		return subtotal
	case DiscountPercent:
		discount := (subtotal * discountValue) / 100
		return subtotal - discount
	case DiscountFixed:
		if discountValue >= subtotal {
			return 0
		}
		return subtotal - discountValue
	default:
		return subtotal
	}
}

func FormatMoney(cents int) string {
	dollars := cents / 100
	centsPart := cents % 100
	return fmt.Sprintf("$%d.%02d", dollars, centsPart)
}

func ValidateReceipt(r *Receipt) []string {
	var errors []string

	if r.PatientID == "" {
		errors = append(errors, "Patient ID is required")
	}

	if r.VisitDate == "" {
		errors = append(errors, "Visit date is required")
	}

	if len(r.LineItems) == 0 {
		errors = append(errors, "At least one line item is required")
	}

	subtotal := 0
	for _, item := range r.LineItems {
		if strings.TrimSpace(item.Description) == "" {
			errors = append(errors, "Line item description is required")
		}
		if item.Quantity <= 0 {
			errors = append(errors, "Line item quantity must be positive")
		}
		if item.UnitPrice < 0 {
			errors = append(errors, "Line item unit price cannot be negative")
		}
		subtotal += item.Quantity * item.UnitPrice
	}

	if r.Subtotal != subtotal {
		errors = append(errors, "Subtotal does not match line items")
	}

	if r.DiscountValue < 0 {
		errors = append(errors, "Discount cannot be negative")
	}
	if r.DiscountType == DiscountPercent && r.DiscountValue > 100 {
		errors = append(errors, "Percentage discount must be between 0 and 100")
	}
	if r.DiscountType == DiscountFixed && r.DiscountValue > subtotal {
		errors = append(errors, "Fixed discount cannot exceed subtotal")
	}

	expectedTotal := CalculateGrandTotal(subtotal, r.DiscountType, r.DiscountValue)
	if r.GrandTotal != expectedTotal {
		errors = append(errors, "Grand total does not match the receipt calculation")
	}
	if r.GrandTotal <= 0 {
		errors = append(errors, "Grand total must be greater than zero")
	}

	return errors
}
