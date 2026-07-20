package models

import (
	"testing"
)

func TestValidateHKIDFormat(t *testing.T) {
	tests := []struct {
		hkid     string
		expected bool
	}{
		{"A123456(7)", true},
		{"AB123456(7)", true},
		{"a123456(7)", true},
		{"A123456(8)", false},
		{"A123456", false},
		{"1234567(8)", false},
		{"ABC12345(6)", false},
	}

	for _, tt := range tests {
		t.Run(tt.hkid, func(t *testing.T) {
			result := ValidateHKIDFormat(tt.hkid)
			if result != tt.expected {
				t.Errorf("ValidateHKIDFormat(%s) = %v, expected %v", tt.hkid, result, tt.expected)
			}
		})
	}
}

func TestCalculateHKIDCheckDigit(t *testing.T) {
	tests := []struct {
		hkid     string
		expected rune
	}{
		{"A123456", '7'},
		{"AB123456", '4'},
	}

	for _, tt := range tests {
		t.Run(tt.hkid, func(t *testing.T) {
			result := CalculateHKIDCheckDigit(tt.hkid)
			if result != tt.expected {
				t.Errorf("CalculateHKIDCheckDigit(%s) = %c, expected %c", tt.hkid, result, tt.expected)
			}
		})
	}
}

func TestHashHKID(t *testing.T) {
	hkid1 := "A123456(7)"
	hkid2 := "a123456(7)"
	hkid3 := "A1234567"

	hash1 := HashHKID(hkid1)
	hash2 := HashHKID(hkid2)
	hash3 := HashHKID(hkid3)

	if hash1 != hash2 {
		t.Error("HKID hashes should be case-insensitive")
	}
	if hash1 != hash3 {
		t.Error("HKID hashes should ignore formatting")
	}
}

func TestMaskHKID(t *testing.T) {
	hkid := "A123456(7)"
	masked := MaskHKID(hkid)
	expected := "A*****6(7)"

	if masked != expected {
		t.Errorf("MaskHKID(%s) = %s, expected %s", hkid, masked, expected)
	}
}

func TestMaskName(t *testing.T) {
	name := "John Doe"
	masked := MaskName(name)
	expected := "J*** D**"

	if masked != expected {
		t.Errorf("MaskName(%s) = %s, expected %s", name, masked, expected)
	}
}