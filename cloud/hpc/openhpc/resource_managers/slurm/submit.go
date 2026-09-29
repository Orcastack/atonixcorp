package slurm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type SubmitOptions struct {
	Partition  string
	Account    string
	TimeLimit  string // e.g. "01:00:00"
	CPUs       int
	MemoryMB   int
	GPUs       int
	JobName    string
	ScriptPath string
	Env        map[string]string
	WorkingDir string
}

type SubmitResult struct {
	JobID string
	Raw   string
}

func (c *SlurmClient) SubmitJob(ctx context.Context, opts SubmitOptions) (*SubmitResult, error) {
	args := []string{}

	if opts.Partition != "" {
		args = append(args, "--partition", opts.Partition)
	}
	if opts.Account != "" {
		args = append(args, "--account", opts.Account)
	}
	if opts.TimeLimit != "" {
		args = append(args, "--time", opts.TimeLimit)
	}
	if opts.CPUs > 0 {
		args = append(args, "--cpus-per-task", fmt.Sprintf("%d", opts.CPUs))
	}
	if opts.MemoryMB > 0 {
		args = append(args, "--mem", fmt.Sprintf("%d", opts.MemoryMB))
	}
	if opts.GPUs > 0 {
		args = append(args, "--gres", fmt.Sprintf("gpu:%d", opts.GPUs))
	}
	if opts.JobName != "" {
		args = append(args, "--job-name", opts.JobName)
	}
	if opts.WorkingDir != "" {
		args = append(args, "--chdir", opts.WorkingDir)
	}

	script := opts.ScriptPath
	if !filepath.IsAbs(script) {
		wd, _ := os.Getwd()
		script = filepath.Join(wd, script)
	}

	args = append(args, script)

	raw, err := c.run(ctx, c.sbatchPath, args...)
	if err != nil {
		return nil, err
	}

	jobID, err := ParseSubmitOutput(raw)
	if err != nil {
		return nil, err
	}

	return &SubmitResult{
		JobID: jobID,
		Raw:   raw,
	}, nil
}
