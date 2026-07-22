package handlers

import "clinic-hcc-app/internal/models"

type Permission = models.Permission

const (
	PermReceiptsCreate    = models.PermReceiptsCreate
	PermReceiptsRead      = models.PermReceiptsRead
	PermReceiptsUpdate    = models.PermReceiptsUpdate
	PermReceiptsFinalize  = models.PermReceiptsFinalize
	PermReceiptsArchive   = models.PermReceiptsArchive
	PermPatientsRead      = models.PermPatientsRead
	PermPatientsCreate    = models.PermPatientsCreate
	PermPatientsUpdate    = models.PermPatientsUpdate
	PermReportsGenerate   = models.PermReportsGenerate
	PermReportsExport     = models.PermReportsExport
	PermSettingsRead      = models.PermSettingsRead
	PermSettingsUpdate    = models.PermSettingsUpdate
	PermBackupManage      = models.PermBackupManage
	PermNotificationsRead = models.PermNotificationsRead
)
