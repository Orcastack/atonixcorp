package datasets

import (
	"context"
	"fmt"
	"io"
	"time"
)

type UploadService struct {
	storage StorageAdapter
	meta    MetadataRepository
}

func NewUploadService(storage StorageAdapter, meta MetadataRepository) *UploadService {
	return &UploadService{
		storage: storage,
		meta:    meta,
	}
}

type UploadRequest struct {
	TenantID    string
	Name        string
	ContentType string
	Labels      map[string]string
	Reader      io.Reader
	SizeBytes   int64
}

func (s *UploadService) Upload(ctx context.Context, req UploadRequest) (*Dataset, error) {
	id := generateDatasetID()
	path := fmt.Sprintf("%s/%s", req.TenantID, id)

	if err := s.storage.Upload(ctx, path, req.Reader, req.SizeBytes, req.ContentType); err != nil {
		return nil, err
	}

	d := &Dataset{
		ID:          id,
		TenantID:    req.TenantID,
		Name:        req.Name,
		SizeBytes:   req.SizeBytes,
		ContentType: req.ContentType,
		StoragePath: path,
		Labels:      req.Labels,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.meta.Save(ctx, d); err != nil {
		return nil, err
	}

	return d, nil
}

func generateDatasetID() string {
	return fmt.Sprintf("ds-%d", time.Now().UnixNano())
}
