package utils

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

type ExecResult struct {
	Stdout string
	Stderr string
}

func RunCommand(ctx context.Context, timeout time.Duration, name string, args ...string) (*ExecResult, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, Wrap(ErrCodeInternal, "command failed", err)
	}

	return &ExecResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}
