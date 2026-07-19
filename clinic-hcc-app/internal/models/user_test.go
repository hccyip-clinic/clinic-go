package models

import (
	"testing"
	"clinic-hcc-app/internal/handlers"
)

func TestDefaultUserHasAllPermissions(t *testing.T) {
	user := DefaultUser()

	allPerms := []handlers.Permission{
		handlers.PermReceiptsCreate,
		handlers.PermReceiptsRead,
		handlers.PermReceiptsUpdate,
		handlers.PermReceiptsFinalize,
		handlers.PermReceiptsArchive,
		handlers.PermPatientsRead,
		handlers.PermPatientsCreate,
		handlers.PermPatientsUpdate,
		handlers.PermReportsGenerate,
		handlers.PermReportsExport,
		handlers.PermSettingsRead,
		handlers.PermSettingsUpdate,
		handlers.PermBackupManage,
		handlers.PermNotificationsRead,
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
		Permissions: []handlers.Permission{handlers.PermReceiptsCreate, handlers.PermPatientsRead},
	}

	if !user.HasPermission(handlers.PermReceiptsCreate) {
		t.Error("Should return true when permission is in list")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionNotAssigned(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []handlers.Permission{handlers.PermReceiptsCreate},
	}

	if user.HasPermission(handlers.PermSettingsUpdate) {
		t.Error("Should return false when permission is not in list")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionsSliceIsEmpty(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []handlers.Permission{},
	}

	if user.HasPermission(handlers.PermReceiptsCreate) {
		t.Error("Should return false when permissions slice is empty")
	}
}

func TestHasPermission_ReturnsFalse_WhenPermissionsSliceIsNil(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: nil,
	}

	if user.HasPermission(handlers.PermReceiptsCreate) {
		t.Error("Should return false when permissions slice is nil")
	}
}

func TestHasPermission_WithEmptyPermissionString(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []handlers.Permission{""},
	}

	if !user.HasPermission("") {
		t.Error("Should match empty string permission if present")
	}
}

func TestHasPermission_WithDuplicatePermissions(t *testing.T) {
	user := &User{
		ID:          "test-user",
		Username:    "test",
		Permissions: []handlers.Permission{handlers.PermReceiptsCreate, handlers.PermReceiptsCreate},
	}

	if !user.HasPermission(handlers.PermReceiptsCreate) {
		t.Error("Should return true even with duplicate permissions")
	}
}

func TestHasPermission_NilReceiver(t *testing.T) {
	var user *User
	if user.HasPermission(handlers.PermReceiptsCreate) {
		t.Error("Should return false when called on nil receiver")
	}
}