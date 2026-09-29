package api

import (
	"net/http"
	"time"

	"atonixcorp/cloud/hpc/api/handlers"
	"atonixcorp/cloud/hpc/api/middleware"
	"atonixcorp/cloud/hpc/jobs"
	"atonixcorp/cloud/hpc/metrics"
	"atonixcorp/cloud/hpc/nodes"
	"atonixcorp/cloud/hpc/openhpc/resource_managers/slurm"
	"atonixcorp/cloud/hpc/scheduler"
	"atonixcorp/cloud/hpc/scheduler/k8s"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Simple backend selector: choose Slurm or K8s based on partition or config.
type BackendSelector struct {
	slurm scheduler.Scheduler
	k8s   scheduler.Scheduler
}

func NewBackendSelector(sl scheduler.Scheduler, ka scheduler.Scheduler) *BackendSelector {
	return &BackendSelector{
		slurm: sl,
		k8s:   ka,
	}
}

func (b *BackendSelector) SelectBackend(tenantID, partition string) (scheduler.Scheduler, string, error) {
	// You can make this smarter (per-tenant, per-partition config).
	if partition == "k8s" {
		return b.k8s, "k8s", nil
	}
	return b.slurm, "slurm", nil
}

func NewRouter(
	jobRepo jobs.Repository,
	eventPublisher jobs.EventPublisher,
	nodeInventory *nodes.InventoryService,
	logRepo jobs.Repository, // same repo if logs stored together
	exporter *metrics.Exporter,
	signingKey []byte,
) http.Handler {
	r := chi.NewRouter()

	// Base chi middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))

	// Security & multi-tenant middleware
	auth := middleware.NewAuthMiddleware(signingKey)
	rate := middleware.NewRateLimiter(200)
	tenantScope := middleware.NewTenantScopeMiddleware()

	r.Use(auth.Middleware)
	r.Use(rate.Middleware)
	r.Use(tenantScope.Middleware)

	// Scheduler backends
	slurmAdapter := slurm.NewSlurmAdapter()
	k8sAdapter, _ := k8s.NewK8sAdapter("hpc-jobs") // namespace
	selector := NewBackendSelector(slurmAdapter, k8sAdapter)

	// Core services
	jobService := jobs.NewService(jobRepo, eventPublisher, selector)
	logService := jobs.NewLogService(logRepo)

	// Handlers
	jobSubmitHandler := handlers.NewJobSubmitHandler(jobService)
	jobStatusHandler := handlers.NewJobStatusHandler(jobService)
	jobLogsHandler := handlers.NewJobLogsHandler(logService)
	nodeListHandler := handlers.NewNodeListHandler(nodeInventory)
	partitionListHandler := handlers.NewPartitionListHandler(nodeInventory)
	metricsHandler := handlers.NewMetricsHandler(exporter)

	// Routes
	r.Route("/hpc", func(h chi.Router) {

		// Jobs
		h.Post("/jobs", jobSubmitHandler.ServeHTTP)
		h.Get("/jobs/{id}", jobStatusHandler.ServeHTTP)
		h.Get("/jobs/{id}/logs", jobLogsHandler.ServeHTTP)

		// Nodes & partitions
		h.Get("/nodes", nodeListHandler.ServeHTTP)
		h.Get("/partitions", partitionListHandler.ServeHTTP)
	})

	// Metrics (Prometheus)
	r.Handle("/metrics", metricsHandler.ServeHTTP)

	// Health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	return r
}
