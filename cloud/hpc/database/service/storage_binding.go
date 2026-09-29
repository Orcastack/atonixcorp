package service

import (
	"atonixcorp/cloud/hpc/database/model"
	"fmt"
)

type StorageBinding struct {
	bindings map[string]*model.StorageBackend
}

func NewStorageBinding() *StorageBinding {
	return &StorageBinding{
		bindings: make(map[string]*model.StorageBackend),
	}
}

func (sb *StorageBinding) Bind(datasetID string, backend *model.StorageBackend) {
	sb.bindings[datasetID] = backend
	fmt.Printf("Storage backend bound: dataset=%s backend=%s\n",
		datasetID, backend.Type)
}

func (sb *StorageBinding) Get(datasetID string) *model.StorageBackend {
	return sb.bindings[datasetID]
}
