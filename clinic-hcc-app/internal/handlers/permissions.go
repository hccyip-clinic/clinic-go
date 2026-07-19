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