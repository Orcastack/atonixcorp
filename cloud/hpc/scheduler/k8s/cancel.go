package k8s

import (
	"context"
)

func (a *K8sAdapter) Cancel(ctx context.Context, jobName string) error {
	return a.client.DeleteJob(ctx, jobName)
}
