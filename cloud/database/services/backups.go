package services

import (
	"fmt"
	"time"

	"atonixcorp/cloud/database/core"
	"atonixcorp/cloud/database/internal"
)

type BackupService struct {
	store core.Store
}

func NewBackupService(store core.Store) *BackupService {
	return &BackupService{store}
}

// TriggerBackup creates backup metadata and triggers backup workflow.
func (b *BackupService) TriggerBackup(instanceID string, backupType string) (*core.Backup, error) {
	internal.Info("Triggering backup for instance %s", instanceID)

	inst, err := b.store.GetInstance(instanceID)
	if err != nil {
		return nil, internal.Wrap("instance lookup failed", err)
	}

	backup := &core.Backup{
		ID:         internal.GenerateID("backup"),
		InstanceID: inst.ID,
		Type:       backupType,
		Location:   fmt.Sprintf("s3://atonix-backups/%s/%s", inst.TenantID, inst.ID),
		CreatedAt:  time.Now(),
	}

	// In real implementation:
	// - Trigger backup job
	// - Upload to MinIO/Ceph S3
	// - Update backup metadata table

	internal.Info("Backup %s created for instance %s", backup.ID, inst.ID)
	return backup, nil
}
