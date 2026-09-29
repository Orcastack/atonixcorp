package slurm

import (
	"errors"
	"strings"
)

func ParseSubmitOutput(raw string) (string, error) {
	// Typical: "Submitted batch job 12345"
	fields := strings.Fields(raw)
	if len(fields) < 4 {
		return "", errors.New("unexpected sbatch output: " + raw)
	}
	return fields[len(fields)-1], nil
}

func ParseStatusOutput(raw string) (*JobStatus, error) {
	line := strings.TrimSpace(raw)
	if line == "" {
		return nil, errors.New("empty slurm status output")
	}

	// Expect: JobID JobName Partition User State Nodes
	fields := strings.Fields(line)
	if len(fields) < 6 {
		return nil, errors.New("unexpected slurm status output: " + raw)
	}

	state := JobState(fields[4])
	if !isValidState(state) {
		state = JobStateUnknown
	}

	return &JobStatus{
		JobID:     fields[0],
		Name:      fields[1],
		Partition: fields[2],
		User:      fields[3],
		State:     state,
		Nodes:     fields[5],
		Raw:       raw,
	}, nil
}

func isValidState(s JobState) bool {
	switch s {
	case JobStatePending,
		JobStateRunning,
		JobStateComplete,
		JobStateFailed,
		JobStateCanceled:
		return true
	default:
		return false
	}
}
