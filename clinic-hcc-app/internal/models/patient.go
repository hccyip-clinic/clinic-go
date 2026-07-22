package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type Patient struct {
	ID        string
	Name      string
	HKID      string
	HKIDHash  string
	Gender    string
	CreatedAt string
	UpdatedAt string
}

func ValidateHKIDFormat(hkid string) bool {
	_, err := NormalizeHKID(hkid)
	return err == nil
}

func CalculateHKIDCheckDigit(hkid string) rune {
	clean := strings.ToUpper(strings.NewReplacer("(", "", ")", "", " ", "", "-", "").Replace(hkid))
	if len(clean) != 7 && len(clean) != 8 {
		return '0'
	}

	if len(clean) == 7 {
		clean = "0" + clean
	}

	weights := []int{9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i, value := range clean {
		var numeric int
		if value >= 'A' && value <= 'Z' {
			numeric = int(value-'A') + 10
		} else if value >= '0' && value <= '9' {
			numeric = int(value - '0')
		} else {
			return '0'
		}
		sum += numeric * weights[i]
	}

	checkDigit := 11 - (sum % 11)
	switch checkDigit {
	case 10:
		return 'A'
	case 11:
		return '0'
	default:
		return rune('0' + checkDigit)
	}
}

func HashHKID(hkid string) string {
	canonical := strings.ToUpper(strings.NewReplacer("(", "", ")", "", " ", "", "-", "").Replace(hkid))
	hash := sha256.Sum256([]byte(canonical))
	return hex.EncodeToString(hash[:])
}

func NormalizeHKID(hkid string) (string, error) {
	clean := strings.ToUpper(strings.NewReplacer("(", "", ")", "", " ", "", "-", "").Replace(strings.TrimSpace(hkid)))
	pattern := regexp.MustCompile(`^([A-Z]{1,2})([0-9]{6})([0-9A])$`)
	matches := pattern.FindStringSubmatch(clean)
	if matches == nil {
		return "", fmt.Errorf("invalid HKID format")
	}

	checkDigit := CalculateHKIDCheckDigit(matches[1] + matches[2])
	if rune(matches[3][0]) != checkDigit {
		return "", fmt.Errorf("invalid HKID check digit")
	}

	return fmt.Sprintf("%s%s(%s)", matches[1], matches[2], matches[3]), nil
}

func MaskHKID(hkid string) string {
	canonical, err := NormalizeHKID(hkid)
	if err != nil {
		return hkid
	}
	return string(canonical[0]) + "*****" + string(canonical[len(canonical)-4]) + canonical[len(canonical)-3:]
}

func MaskName(name string) string {
	words := strings.Fields(name)
	masked := make([]string, len(words))

	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			masked[i] = string(runes[0]) + strings.Repeat("*", len(runes)-1)
		}
	}

	return strings.Join(masked, " ")
}
