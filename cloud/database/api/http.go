package api

import (
	"net/http"

	"atonixcorp/cloud/database/core"
	"atonixcorp/cloud/database/services"

	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	provisioning *services.ProvisioningService
	backups      *services.BackupService
	metrics      *services.MetricsService
}

func NewAPIServer(
	controller *core.Controller,
	store core.Store,
	os core.OpenStackClient,
) *APIServer {
	return &APIServer{
		provisioning: services.NewProvisioningService(controller),
		backups:      services.NewBackupService(store),
		metrics:      services.NewMetricsService(store, os),
	}
}

func (s *APIServer) Router() http.Handler {
	r := chi.NewRouter()

	// Instances
	r.Post("/tenants/{tenant_id}/instances", s.CreateInstance)
	r.Get("/tenants/{tenant_id}/instances/{instance_id}", s.GetInstance)

	// Backups
	r.Post("/tenants/{tenant_id}/instances/{instance_id}/backups", s.TriggerBackup)

	// Metrics
	r.Get("/tenants/{tenant_id}/instances/{instance_id}/metrics", s.GetMetrics)

	return r
}
