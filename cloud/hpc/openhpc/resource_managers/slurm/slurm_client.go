package slurm

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type SlurmClient struct {
	sbatchPath  string
	squeuePath  string
	sacctPath   string
	scancelPath string
	timeout     time.Duration
}

func NewSlurmClient() *SlurmClient {
	return &SlurmClient{
		sbatchPath:  "sbatch",
		squeuePath:  "squeue",
		sacctPath:   "sacct",
		scancelPath: "scancel",
		timeout:     30 * time.Second,
	}
}

func (c *SlurmClient) run(ctx context.Context, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", &CommandError{
			Cmd:    name,
			Args:   args,
			Stderr: stderr.String(),
			Err:    err,
		}
	}

	return stdout.String(), nil
}

type CommandError struct {
	Cmd    string
	Args   []string
	Stderr string
	Err    error
}

func (e *CommandError) Error() string {
	return "slurm command failed: " + e.Cmd + " " + e.Stderr
}
