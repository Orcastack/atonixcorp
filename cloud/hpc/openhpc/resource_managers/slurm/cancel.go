package slurm

import "context"

func (c *SlurmClient) CancelJob(ctx context.Context, jobID string) error {
	_, err := c.run(ctx, c.scancelPath, jobID)
	return err
}
