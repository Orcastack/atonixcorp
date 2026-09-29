package nodes

import (
	"context"
	"time"
)

type PartitionService struct {
	inventory *InventoryService
}

func NewPartitionService(inv *InventoryService) *PartitionService {
	return &PartitionService{inventory: inv}
}

func (s *PartitionService) CreatePartition(ctx context.Context, p Partition) error {
	now := time.Now().UTC()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	p.UpdatedAt = now
	if p.Labels == nil {
		p.Labels = make(map[string]string)
	}
	p.Active = true
	return s.inventory.repo.SavePartition(ctx, p)
}

func (s *PartitionService) AssignNodeToPartition(ctx context.Context, nodeID, partitionName string) error {
	node, err := s.inventory.repo.GetNode(ctx, nodeID)
	if err != nil {
		return err
	}

	p, err := s.inventory.repo.GetPartition(ctx, partitionName)
	if err != nil {
		return err
	}

	// add partition to node
	if !containsString(node.Partitions, partitionName) {
		node.Partitions = append(node.Partitions, partitionName)
		node.UpdatedAt = time.Now().UTC()
		if err := s.inventory.repo.SaveNode(ctx, node); err != nil {
			return err
		}
	}

	// add node to partition
	if !containsString(p.Nodes, nodeID) {
		p.Nodes = append(p.Nodes, nodeID)
		p.UpdatedAt = time.Now().UTC()
		if err := s.inventory.repo.SavePartition(ctx, p); err != nil {
			return err
		}
	}

	return nil
}

func (s *PartitionService) ListPartitionNodes(ctx context.Context, partitionName string) ([]Node, error) {
	p, err := s.inventory.repo.GetPartition(ctx, partitionName)
	if err != nil {
		return nil, err
	}

	allNodes, err := s.inventory.ListAllNodes(ctx)
	if err != nil {
		return nil, err
	}

	var result []Node
	for _, n := range allNodes {
		if containsString(p.Nodes, n.ID) {
			result = append(result, n)
		}
	}
	return result, nil
}

func containsString(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}
