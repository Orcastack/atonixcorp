package jobs

import (
	"context"
	"time"
)

type LogService struct {
	repo Repository
}

func NewLogService(repo Repository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) Append(ctx context.Context, jobID, source, content string) error {
	log := &JobLog{
		JobID:     jobID,
		Source:    source,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.SaveLog(ctx, log)
}

func (s *LogService) List(ctx context.Context, jobID string) ([]*JobLog, error) {
	return s.repo.ListLogs(ctx, jobID)
}
