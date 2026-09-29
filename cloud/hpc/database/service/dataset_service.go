package service

import (
	"atonixcorp/cloud/hpc/database/model"
	"atonixcorp/cloud/hpc/database/registry"
	"fmt"
)

type DatasetService struct {
	DB *DatabaseService
}

func NewDatasetService(db *DatabaseService) *DatasetService {
	return &DatasetService{DB: db}
}

func (ds *DatasetService) Create(meta *model.Dataset) {
	ds.DB.RegisterDataset(meta)
}

func (ds *DatasetService) Get(id string) *registry.DatasetMetadata {
	return ds.DB.Registry.Get(id)
}

func (ds *DatasetService) Delete(id string) {
	ds.DB.Registry.Delete(id)
}

func (ds *DatasetService) AddVersion(id string, notes string, timestamp string) {
	versions := ds.DB.Versions.GetVersions(id)
	next := len(versions) + 1

	ds.DB.Versions.AddVersion(id, registry.Version{
		Number:    next,
		Timestamp: timestamp,
		Notes:     notes,
	})
}

func (ds *DatasetService) AddLineage(id string, parent string) {
	fmt.Printf("Lineage added: %s <- %s\n", id, parent)
}
