package k8s

import (
	"context"
	"fmt"

	"atonixcorp/cloud/hpc/scheduler"
)

type SubmitOptions struct {
	Image   string
	Command []string
	Args    []string
}

type SubmitResult struct {
	JobName string
}

type K8sAdapter struct {
	client *K8sClient
}

func NewK8sAdapter(namespace string) (*K8sAdapter, error) {
	client, err := NewK8sClient(namespace)
	if err != nil {
		return nil, err
	}
	return &K8sAdapter{client: client}, nil
}

func (a *K8sAdapter) Submit(ctx context.Context, spec scheduler.JobSpec) (string, error) {
	template := JobTemplate{
		Name:       spec.JobName,
		Image:      spec.Env["IMAGE"], // or from spec directly if you add it
		Command:    []string{"/bin/sh", "-c"},
		Args:       []string{fmt.Sprintf("bash %s", spec.ScriptPath)},
		CPUs:       fmt.Sprintf("%d", spec.CPUs),
		Memory:     fmt.Sprintf("%dMi", spec.MemoryMB),
		GPUs:       fmt.Sprintf("%d", spec.GPUs),
		Env:        spec.Env,
		WorkingDir: spec.WorkingDir,
		Labels: map[string]string{
			"partition": spec.Partition,
			"account":   spec.Account,
		},
	}

	job := BuildJob(template)

	_, err := a.client.CreateJob(ctx, job)
	if err != nil {
		return "", err
	}

	return job.Name, nil
}
