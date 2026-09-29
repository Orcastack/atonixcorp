package model

// StorageBackend represents the underlying storage system
// used for a dataset (Lustre, BeeGFS, GPFS, S3, Ceph, etc.)
type StorageBackend struct {
	ID          string
	Type        string // lustre, beegfs, gpfs, s3, ceph, local
	MountPoint  string
	Endpoint    string
	Description string
	CreatedAt   string
}
