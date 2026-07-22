package models

import "testing"

func TestDefaultUserHasAllPermissions(t *testing.T) {
	user := DefaultUser()

	allPerms := []Permission{
		PermReceiptsCreate, PermReceiptsRead, PermReceiptsUpdate,
		PermReceiptsFinalize, PermReceiptsArchive, PermPatientsRead,
		PermPatientsCreate, PermPatientsUpdate, PermReportsGenerate,
		PermReportsExport, PermSettingsRead, PermSettingsUpdate,
		PermBackupManage, PermNotificationsRead,
	}

	for _, perm := range allPerms {
		t.Run(string(perm), func(t *testing.T) {
			if !user.HasPermission(perm) {
				t.Errorf("Default user should have permission %s", perm)
			}
		})
	}
}

func TestHasPermission_ReturnsTrue_WhenPermissionExists(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []Permission{PermReceiptsCreate, PermPatientsRead},
	}

	if !user.HasPermission(PermReceiptsCreate) {
		t.Error("Should return true when permission is in list")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionNotAssigned(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []Permission{PermReceiptsCreate},
	}

	if user.HasPermission(PermSettingsUpdate) {
		t.Error("Should return false when permission is not in list")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionsSliceIsEmpty(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []Permission{},
	}

	if user.HasPermission(PermReceiptsCreate) {
		t.Error("Should return false when permissions slice is empty")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionsSliceIsNil(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: nil,
	}

	if user.HasPermission(PermReceiptsCreate) {
		t.Error("Should return false when permissions slice is nil")
	}
}

func TestHasPermission_WithEmptyPermissionString(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []Permission{""},
	}

	if !user.HasPermission("") {
		t.Error("Should match empty string permission if present")
	}
}

func TestHasPermission_WithDuplicatePermissions(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []Permission{PermReceiptsCreate, PermReceiptsCreate},
	}

	if !user.HasPermission(PermReceiptsCreate) {
		t.Error("Should return true even with duplicate permissions")
	}
}

func TestHasPermission_NilReceiver(t *testing.T) {
	var user *User
	if user.HasPermission(PermReceiptsCreate) {
		t.Error("Should return false when called on nil receiver")
	}
}
