package jobs

import (
	"context"
)

type Repository interface {
	SaveJob(ctx context.Context, job *Job) error
	GetJob(ctx context.Context, id string) (*Job, error)
	ListJobsByTenant(ctx context.Context, tenantID string) ([]*Job, error)
	SaveLog(ctx context.Context, log *JobLog) error
	ListLogs(ctx context.Context, jobID string) ([]*JobLog, error)
}
