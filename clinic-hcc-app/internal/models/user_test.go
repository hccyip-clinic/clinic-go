package models

import (
	"testing"
	"clinic-hcc-app/internal/handlers"
)

func TestDefaultUserHasAllPermissions(t *testing.T) {
	user := DefaultUser()

	tests := []struct {
		perm handlers.Permission
	}{
		{handlers.PermReceiptsCreate},
		{handlers.PermReceiptsRead},
		{handlers.PermPatientsRead},
		{handlers.PermSettingsUpdate},
	}

	for _, tt := range tests {
		t.Run(string(tt.perm), func(t *testing.T) {
			if !user.HasPermission(tt.perm) {
				t.Errorf("Default user should have permission %s", tt.perm)
			}
		})
	}
}