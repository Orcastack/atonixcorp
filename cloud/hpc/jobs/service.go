package jobs

import (
	"context"
	"time"

	"atonixcorp/cloud/hpc/scheduler"
)

type BackendSelector interface {
	SelectBackend(tenantID, partition string) (scheduler.Scheduler, string, error)
}

type Service struct {
	repo     Repository
	events   EventPublisher
	selector BackendSelector
}

func NewService(repo Repository, events EventPublisher, selector BackendSelector) *Service {
	return &Service{
		repo:     repo,
		events:   events,
		selector: selector,
	}
}

func (s *Service) CreateAndSubmit(ctx context.Context, req CreateJobRequest) (*Job, error) {
	if err := ValidateCreateJob(req); err != nil {
		return nil, err
	}

	backend, backendName, err := s.selector.SelectBackend(req.TenantID, req.Partition)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	job := &Job{
		ID:         generateJobID(), // implement ID generator
		TenantID:   req.TenantID,
		UserID:     req.UserID,
		Name:       req.Name,
		Backend:    backendName,
		Partition:  req.Partition,
		Account:    req.Account,
		ScriptPath: req.ScriptPath,
		CPUs:       req.CPUs,
		MemoryMB:   req.MemoryMB,
		GPUs:       req.GPUs,
		WorkingDir: req.WorkingDir,
		Env:        req.Env,
		State:      JobStatePending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}

	_ = s.events.Publish(Event{
		Type:      EventJobCreated,
		JobID:     job.ID,
		TenantID:  job.TenantID,
		Timestamp: now,
		Payload:   map[string]any{"name": job.Name},
	})

	spec := scheduler.JobSpec{
		JobName:    job.ID,
		ScriptPath: job.ScriptPath,
		Partition:  job.Partition,
		Account:    job.Account,
		TimeLimit:  "", // extend later
		CPUs:       job.CPUs,
		MemoryMB:   job.MemoryMB,
		GPUs:       job.GPUs,
		WorkingDir: job.WorkingDir,
		Env:        job.Env,
	}

	schedulerID, err := backend.Submit(ctx, spec)
	if err != nil {
		job.State = JobStateFailed
		job.UpdatedAt = time.Now().UTC()
		_ = s.repo.SaveJob(ctx, job)
		return nil, err
	}

	job.SchedulerID = schedulerID
	job.State = JobStateQueued
	job.UpdatedAt = time.Now().UTC()
	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}

	_ = s.events.Publish(Event{
		Type:      EventJobQueued,
		JobID:     job.ID,
		TenantID:  job.TenantID,
		Timestamp: time.Now().UTC(),
		Payload:   map[string]any{"scheduler_id": schedulerID},
	})

	return job, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Job, error) {
	return s.repo.GetJob(ctx, id)
}

func (s *Service) RefreshStatus(ctx context.Context, id string) (*Job, error) {
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}

	backend, _, err := s.selector.SelectBackend(job.TenantID, job.Partition)
	if err != nil {
		return nil, err
	}

	status, err := backend.Status(ctx, job.SchedulerID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	job.UpdatedAt = now

	switch status.State {
	case "PENDING":
		job.State = JobStateQueued
	case "RUNNING":
		job.State = JobStateRunning
		if job.StartedAt == nil {
			job.StartedAt = &now
		}
	case "COMPLETED":
		job.State = JobStateSucceeded
		if job.FinishedAt == nil {
			job.FinishedAt = &now
		}
	case "FAILED":
		job.State = JobStateFailed
		if job.FinishedAt == nil {
			job.FinishedAt = &now
		}
	case "CANCELLED":
		job.State = JobStateCanceled
		if job.FinishedAt == nil {
			job.FinishedAt = &now
		}
	}

	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}

	return job, nil
}

func (s *Service) Cancel(ctx context.Context, id string) (*Job, error) {
	job, err := s.repo.GetJob(ctx, id)
	if err != nil {
		return nil, err
	}

	backend, _, err := s.selector.SelectBackend(job.TenantID, job.Partition)
	if err != nil {
		return nil, err
	}

	if err := backend.Cancel(ctx, job.SchedulerID); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	job.State = JobStateCanceled
	job.UpdatedAt = now
	job.FinishedAt = &now

	if err := s.repo.SaveJob(ctx, job); err != nil {
		return nil, err
	}

	_ = s.events.Publish(Event{
		Type:      EventJobCanceled,
		JobID:     job.ID,
		TenantID:  job.TenantID,
		Timestamp: now,
		Payload:   map[string]any{},
	})

	return job, nil
}

// TODO: replace with your own ID generator (UUID, ULID, etc.)
func generateJobID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
