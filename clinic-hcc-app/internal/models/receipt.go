package models

import (
	"fmt"
	"math/rand"
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

func GenerateReceiptNumber(prefix string) string {
	now := time.Now()
	date := now.Format("20060102")

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	random := make([]byte, 5)
	for i := range random {
		random[i] = chars[rand.Intn(len(chars))]
	}

	return fmt.Sprintf("%s-%s-%s", prefix, date, string(random))
}

func CalculateSubtotal(items []LineItem) int {
	total := 0
	for _, item := range items {
		total += item.Subtotal
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

	if r.GrandTotal <= 0 {
		errors = append(errors, "Grand total must be greater than zero")
	}

	return errors
}