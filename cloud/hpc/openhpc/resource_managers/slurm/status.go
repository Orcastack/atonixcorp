package slurm

import (
	"context"
)

type JobState string

const (
	JobStatePending  JobState = "PENDING"
	JobStateRunning  JobState = "RUNNING"
	JobStateComplete JobState = "COMPLETED"
	JobStateFailed   JobState = "FAILED"
	JobStateCanceled JobState = "CANCELLED"
	JobStateUnknown  JobState = "UNKNOWN"
)

type JobStatus struct {
	JobID     string
	Name      string
	State     JobState
	Partition string
	User      string
	Nodes     string
	Raw       string
}

func (c *SlurmClient) GetJobStatus(ctx context.Context, jobID string) (*JobStatus, error) {
	args := []string{
		"--job", jobID,
		"--noheader",
		"--format", "JobID,JobName,Partition,UserName,State,Nodes",
	}

	raw, err := c.run(ctx, c.squeuePath, args...)
	if err != nil {
		// if not in squeue, try sacct for completed/failed
		history, herr := c.run(ctx, c.sacctPath,
			"--jobs", jobID,
			"--noheader",
			"--format", "JobID,JobName,Partition,User,State,NodeList",
		)
		if herr != nil {
			return nil, err
		}
		status, perr := ParseStatusOutput(history)
		if perr != nil {
			return nil, perr
		}
		return status, nil
	}

	status, err := ParseStatusOutput(raw)
	if err != nil {
		return nil, err
	}
	return status, nil
}
