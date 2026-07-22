package models

import (
	"testing"
)

func TestCalculateGrandTotal(t *testing.T) {
	tests := []struct {
		subtotal      int
		discountType  DiscountType
		discountValue int
		expected      int
	}{
		{10000, DiscountNone, 0, 10000},
		{10000, DiscountPercent, 10, 9000},
		{10000, DiscountPercent, 25, 7500},
		{10000, DiscountFixed, 2000, 8000},
		{10000, DiscountFixed, 15000, 0},
	}

	for _, tt := range tests {
		t.Run(string(tt.discountType), func(t *testing.T) {
			result := CalculateGrandTotal(tt.subtotal, tt.discountType, tt.discountValue)
			if result != tt.expected {
				t.Errorf("CalculateGrandTotal(%d, %s, %d) = %d, expected %d",
					tt.subtotal, tt.discountType, tt.discountValue, result, tt.expected)
			}
		})
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		cents    int
		expected string
	}{
		{0, "$0.00"},
		{100, "$1.00"},
		{12345, "$123.45"},
		{99, "$0.99"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := FormatMoney(tt.cents)
			if result != tt.expected {
				t.Errorf("FormatMoney(%d) = %s, expected %s", tt.cents, result, tt.expected)
			}
		})
	}
}

func TestValidateReceipt(t *testing.T) {
	validReceipt := &Receipt{
		PatientID: "patient-123",
		VisitDate: "2026-07-20",
		LineItems: []LineItem{
			{Description: "Consultation", Quantity: 1, UnitPrice: 50000, Subtotal: 50000},
		},
		Subtotal:     50000,
		DiscountType: DiscountNone,
		GrandTotal:   50000,
		Status:       StatusDraft,
	}

	errors := ValidateReceipt(validReceipt)
	if len(errors) != 0 {
		t.Errorf("Valid receipt should have no errors, got %v", errors)
	}

	invalidReceipt := &Receipt{
		PatientID:  "",
		VisitDate:  "",
		LineItems:  []LineItem{},
		GrandTotal: 0,
	}

	errors = ValidateReceipt(invalidReceipt)
	if len(errors) == 0 {
		t.Error("Invalid receipt should have errors")
	}
}

func TestGenerateReceiptNumber(t *testing.T) {
	number, err := GenerateReceiptNumber("RCP")
	if err != nil {
		t.Fatalf("GenerateReceiptNumber failed: %v", err)
	}

	if len(number) != 19 {
		t.Errorf("Receipt number length should be 19, got %d", len(number))
	}
}
