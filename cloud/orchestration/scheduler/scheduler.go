package scheduler

type Scheduler struct{}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) PickNode(cpu int, ram int) (string, error) {
	// TODO: query your compute nodes
	// TODO: pick the best node based on free resources
	return "node-1", nil
}
