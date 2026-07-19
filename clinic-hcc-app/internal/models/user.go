package models

import "clinic-hcc-app/internal/handlers"

type User struct {
	ID          string
	Username    string
	Permissions []handlers.Permission
}

func (u *User) HasPermission(perm handlers.Permission) bool {
	if u == nil {
		return false
	}
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