package datasets

import (
	"context"
	"time"
)

type Dataset struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	Name        string            `json:"name"`
	SizeBytes   int64             `json:"size_bytes"`
	ContentType string            `json:"content_type"`
	StoragePath string            `json:"storage_path"`
	Labels      map[string]string `json:"labels"`
	CreatedAt   time.Time         `json:"created_at"`
}

type MetadataRepository interface {
	Save(ctx context.Context, d *Dataset) error
	Get(ctx context.Context, id string) (*Dataset, error)
	ListByTenant(ctx context.Context, tenantID string) ([]*Dataset, error)
}
