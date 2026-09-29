package datasets

import (
	"context"
	"io"
)

type DownloadService struct {
	storage StorageAdapter
	meta    MetadataRepository
}

func NewDownloadService(storage StorageAdapter, meta MetadataRepository) *DownloadService {
	return &DownloadService{
		storage: storage,
		meta:    meta,
	}
}

func (s *DownloadService) Download(ctx context.Context, datasetID string) (io.ReadCloser, *Dataset, error) {
	d, err := s.meta.Get(ctx, datasetID)
	if err != nil {
		return nil, nil, err
	}

	reader, err := s.storage.Download(ctx, d.StoragePath)
	if err != nil {
		return nil, nil, err
	}

	return reader, d, nil
}
