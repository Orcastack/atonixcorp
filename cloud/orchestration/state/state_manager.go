package state

import "fmt"

type StateManager struct{}

func NewStateManager() *StateManager {
	return &StateManager{}
}

// Desired state: what the user wants
func (s *StateManager) SetDesired(id string, state string) error {
	fmt.Println("[STATE] desired:", id, "→", state)
	// TODO: write to DB
	return nil
}

// Actual state: what the system is currently doing
func (s *StateManager) SetActual(id string, state string) error {
	fmt.Println("[STATE] actual:", id, "→", state)
	// TODO: write to DB
	return nil
}

// Node placement
func (s *StateManager) SetActualNode(id string, node string) error {
	fmt.Println("[STATE] instance:", id, "placed on node:", node)
	// TODO: write to DB
	return nil
}
