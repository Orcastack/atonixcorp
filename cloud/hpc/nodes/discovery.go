package nodes

import (
	"context"
	"time"
)

type DiscoverySource string

const (
	DiscoverySourceSlurm DiscoverySource = "slurm"
	DiscoverySourceK8s   DiscoverySource = "k8s"
	DiscoverySourceAgent DiscoverySource = "agent"
)

type DiscoveryClient interface {
	DiscoverNodes(ctx context.Context) ([]Node, error)
}

type DiscoveryService struct {
	inventory *InventoryService
	client    DiscoveryClient
}

func NewDiscoveryService(inv *InventoryService, client DiscoveryClient) *DiscoveryService {
	return &DiscoveryService{
		inventory: inv,
		client:    client,
	}
}

func (s *DiscoveryService) SyncClusterNodes(ctx context.Context) error {
	nodes, err := s.client.DiscoverNodes(ctx)
	if err != nil {
		return err
	}

	for _, n := range nodes {
		_ = s.inventory.RegisterNode(ctx, n)
	}
	return nil
}

// Example periodic sync (run from main or a supervisor)
func (s *DiscoveryService) StartPeriodicSync(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = s.SyncClusterNodes(ctx)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
