package lifecycle

import (
	"fmt"

	"atonixcorp/cloud/compute/rpc"
)

type Manager struct {
	RPC rpc.RPC
}

func NewManager(backend rpc.RPC) *Manager {
	return &Manager{RPC: backend}
}

// High-level API

func (m *Manager) Start(vm VM) error {
	fmt.Println("Lifecycle: start VM", vm.ID)
	return m.sendToNode(vm.NodeID, "vm.start", vm.ID)
}

func (m *Manager) Stop(vm VM) error {
	fmt.Println("Lifecycle: stop VM", vm.ID)
	return m.sendToNode(vm.NodeID, "vm.stop", vm.ID)
}

func (m *Manager) Reboot(vm VM) error {
	fmt.Println("Lifecycle: reboot VM", vm.ID)
	return m.sendToNode(vm.NodeID, "vm.reboot", vm.ID)
}

func (m *Manager) Pause(vm VM) error {
	fmt.Println("Lifecycle: pause VM", vm.ID)
	return m.sendToNode(vm.NodeID, "vm.pause", vm.ID)
}

func (m *Manager) Migrate(vm VM) error {
	fmt.Println("Lifecycle: migrate VM", vm.ID, "to", vm.TargetNodeID)

	// 1. stop on source
	if err := m.sendToNode(vm.NodeID, "vm.stop", vm.ID); err != nil {
		return err
	}

	// 2. start on target
	if err := m.sendToNode(vm.TargetNodeID, "vm.start", vm.ID); err != nil {
		return err
	}

	return nil
}

// Internal helper

func (m *Manager) sendToNode(nodeID, msgType string, payload any) error {
	return m.RPC.Send(rpc.RPCMessage{
		Type:    msgType,
		Payload: payload,
	})
}
