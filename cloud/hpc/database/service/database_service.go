package service

import (
	"fmt"

	"atonixcorp/cloud/hpc/database/io"
	"atonixcorp/cloud/hpc/database/model"
	"atonixcorp/cloud/hpc/database/monitoring"
	"atonixcorp/cloud/hpc/database/registry"
)

type DatabaseService struct {
	Registry  *registry.DatasetRegistry
	Versions  *registry.VersionManager
	Access    *registry.AccessControl
	IO        *io.ParallelIO
	MPIIO     *io.MPIIO
	Cache     *io.CacheEngine
	Mounts    *io.Mounts
	Metrics   map[string]*monitoring.DatasetMetrics
	IOMetrics map[string]*monitoring.IOMetrics
	Health    map[string]*monitoring.HealthStatus
}

func NewDatabaseService() *DatabaseService {
	return &DatabaseService{
		Registry:  registry.NewDatasetRegistry(),
		Versions:  registry.NewVersionManager(),
		Access:    registry.NewAccessControl(),
		IO:        io.NewParallelIO(),
		MPIIO:     io.NewMPIIO(),
		Cache:     io.NewCacheEngine(),
		Mounts:    io.NewMounts(),
		Metrics:   make(map[string]*monitoring.DatasetMetrics),
		IOMetrics: make(map[string]*monitoring.IOMetrics),
		Health:    make(map[string]*monitoring.HealthStatus),
	}
}

func (db *DatabaseService) RegisterDataset(meta *model.Dataset) {
	db.Registry.Register(&registry.DatasetMetadata{
		ID:          meta.ID,
		Name:        meta.Name,
		Format:      meta.Format,
		Owner:       meta.Owner,
		Tags:        meta.Tags,
		Description: meta.Description,
		CreatedAt:   meta.CreatedAt,
		UpdatedAt:   meta.UpdatedAt,
	})

	db.Metrics[meta.ID] = monitoring.NewDatasetMetrics(meta.ID)
	db.IOMetrics[meta.ID] = monitoring.NewIOMetrics(meta.ID)
	db.Health[meta.ID] = monitoring.NewHealthStatus(meta.ID)

	fmt.Printf("Dataset registered in DatabaseService: %s\n", meta.ID)
}
