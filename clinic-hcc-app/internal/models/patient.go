package models

import (
	"crypto/sha256"
	"encoding/hex"
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
	pattern := regexp.MustCompile(`^[A-Z]{1,2}[0-9]{6}\([0-9A]\)$`)
	return pattern.MatchString(strings.ToUpper(hkid))
}

func CalculateHKIDCheckDigit(hkid string) rune {
	clean := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(hkid, "(", ""), ")", ""))
	
	if len(clean) < 7 {
		return '0'
	}

	weights := []int{9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0

	firstChar := int(clean[0] - 'A' + 10)
	sum += firstChar * weights[0]

	for i := 1; i < 7 && i < len(clean); i++ {
		digit := int(clean[i] - '0')
		sum += digit * weights[i]
	}

	remainder := sum % 11
	checkDigit := 11 - remainder

	if checkDigit == 11 {
		return '0'
	} else if checkDigit == 10 {
		return 'A'
	}
	return rune(checkDigit + '0')
}

func HashHKID(hkid string) string {
	normalized := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(hkid, "(", ""), ")", ""))
	hash := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(hash[:])
}

func MaskHKID(hkid string) string {
	if len(hkid) < 9 {
		return hkid
	}
	return string(hkid[0]) + "*****" + hkid[len(hkid)-3:]
}

func MaskName(name string) string {
	words := strings.Fields(name)
	masked := make([]string, len(words))
	
	for i, word := range words {
		if len(word) > 0 {
			masked[i] = string(word[0]) + strings.Repeat("*", len(word)-1)
		}
	}
	
	return strings.Join(masked, " ")
}