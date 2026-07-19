# Task 2: Permission Constants and User Model

**Goal:** Define permission constants and user model with `HasPermission()` method.

**Files:**
- Create: `internal/handlers/permissions.go`
- Create: `internal/models/user.go`
- Test: `internal/handlers/permissions_test.go`

---

## Step 1: Create Permission Constants

**File:** `internal/handlers/permissions.go`

```go
package handlers

type Permission string

const (
	PermReceiptsCreate   Permission = "receipts:create"
	PermReceiptsRead     Permission = "receipts:read"
	PermReceiptsUpdate   Permission = "receipts:update"
	PermReceiptsFinalize Permission = "receipts:finalize"
	PermReceiptsArchive  Permission = "receipts:archive"
	PermPatientsRead     Permission = "patients:read"
	PermPatientsCreate   Permission = "patients:create"
	PermPatientsUpdate   Permission = "patients:update"
	PermReportsGenerate  Permission = "reports:generate"
	PermReportsExport    Permission = "reports:export"
	PermSettingsRead     Permission = "settings:read"
	PermSettingsUpdate   Permission = "settings:update"
	PermBackupManage     Permission = "backup:manage"
	PermNotificationsRead Permission = "notifications:read"
)
```

---

## Step 2: Create User Model

**File:** `internal/models/user.go`

```go
package models

import "clinic-hcc-app/internal/handlers"

type User struct {
	ID          string
	Username    string
	Permissions []handlers.Permission
}

func (u *User) HasPermission(perm handlers.Permission) bool {
	for _, p := range u.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

func DefaultUser() *User {
	return &User{
		ID:       "user-default",
		Username: "practitioner",
		Permissions: []handlers.Permission{
			handlers.PermReceiptsCreate, handlers.PermReceiptsRead,
			handlers.PermReceiptsUpdate, handlers.PermReceiptsFinalize,
			handlers.PermReceiptsArchive,
			handlers.PermPatientsRead, handlers.PermPatientsCreate,
			handlers.PermPatientsUpdate,
			handlers.PermReportsGenerate, handlers.PermReportsExport,
			handlers.PermSettingsRead, handlers.PermSettingsUpdate,
			handlers.PermBackupManage, handlers.PermNotificationsRead,
		},
	}
}
```

---

## Step 3: Create Tests

**File:** `internal/handlers/permissions_test.go`

```go
package handlers

import (
	"testing"
	"clinic-hcc-app/internal/models"
)

func TestDefaultUserHasAllPermissions(t *testing.T) {
	user := models.DefaultUser()

	tests := []struct {
		perm Permission
	}{
		{PermReceiptsCreate},
		{PermReceiptsRead},
		{PermPatientsRead},
		{PermSettingsUpdate},
	}

	for _, tt := range tests {
		t.Run(string(tt.perm), func(t *testing.T) {
			if !user.HasPermission(tt.perm) {
				t.Errorf("Default user should have permission %s", tt.perm)
			}
		})
	}
}
```

---

## Step 4: Run Tests

```bash
cd clinic-hcc-app
go test ./internal/handlers/... ./internal/models/... -v
```

**Expected:** PASS

---

## Step 5: Commit

```bash
git add internal/handlers/permissions.go internal/models/user.go internal/handlers/permissions_test.go
git commit -m "feat: add permission constants and user model"
```