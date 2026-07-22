package models

type Permission string

const (
	PermReceiptsCreate    Permission = "receipts:create"
	PermReceiptsRead      Permission = "receipts:read"
	PermReceiptsUpdate    Permission = "receipts:update"
	PermReceiptsFinalize  Permission = "receipts:finalize"
	PermReceiptsArchive   Permission = "receipts:archive"
	PermPatientsRead      Permission = "patients:read"
	PermPatientsCreate    Permission = "patients:create"
	PermPatientsUpdate    Permission = "patients:update"
	PermReportsGenerate   Permission = "reports:generate"
	PermReportsExport     Permission = "reports:export"
	PermSettingsRead      Permission = "settings:read"
	PermSettingsUpdate    Permission = "settings:update"
	PermBackupManage      Permission = "backup:manage"
	PermNotificationsRead Permission = "notifications:read"
)

type User struct {
	ID          string
	Username    string
	Permissions []Permission
}

func (u *User) HasPermission(perm Permission) bool {
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
		Permissions: []Permission{
			PermReceiptsCreate, PermReceiptsRead,
			PermReceiptsUpdate, PermReceiptsFinalize,
			PermReceiptsArchive,
			PermPatientsRead, PermPatientsCreate,
			PermPatientsUpdate,
			PermReportsGenerate, PermReportsExport,
			PermSettingsRead, PermSettingsUpdate,
			PermBackupManage, PermNotificationsRead,
		},
	}
}
