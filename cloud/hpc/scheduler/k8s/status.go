package k8s

import (
	"context"

	"atonixcorp/cloud/hpc/scheduler"
)

func (a *K8sAdapter) Status(ctx context.Context, jobName string) (*scheduler.JobStatus, error) {
	job, err := a.client.GetJob(ctx, jobName)
	if err != nil {
		return nil, err
	}

	state := "UNKNOWN"
	if job.Status.Active > 0 {
		state = "RUNNING"
	} else if job.Status.Succeeded > 0 {
		state = "COMPLETED"
	} else if job.Status.Failed > 0 {
		state = "FAILED"
	}

	nodes := ""
	pods, err := a.client.ListPodsForJob(ctx, jobName)
	if err == nil && len(pods.Items) > 0 {
		for _, p := range pods.Items {
			if len(p.Spec.NodeName) > 0 {
				nodes = p.Spec.NodeName
				break
			}
		}
	}

	return &scheduler.JobStatus{
		JobID:     job.Name,
		Name:      job.Name,
		State:     state,
		Partition: job.Labels["partition"],
		User:      job.Labels["user"],
		Nodes:     nodes,
		Raw:       "",
	}, nil
}
