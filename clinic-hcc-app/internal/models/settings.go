package models

import "fmt"

type Settings struct {
	ID                     int
	ClinicName             string
	ClinicAddress          string
	ClinicPhone            string
	PractitionerName       string
	PractitionerRegistration string
	ReceiptPrefix          string
	RetentionYears         int
	UpdatedAt              string
}

func ValidateSettings(s *Settings) []string {
	var errors []string

	if s.ClinicName == "" {
		errors = append(errors, "Clinic name is required")
	}

	if s.PractitionerName == "" {
		errors = append(errors, "Practitioner name is required")
	}

	if s.ReceiptPrefix == "" {
		errors = append(errors, "Receipt prefix is required")
	}

	if s.RetentionYears < 1 {
		errors = append(errors, "Retention years must be at least 1")
	}

	return errors
}

func GetFinancialYear(date string) (startYear, endYear int) {
	var year, month, day int
	fmt.Sscanf(date, "%d-%d-%d", &year, &month, &day)

	if month >= 4 {
		return year, year + 1
	}
	return year - 1, year
}