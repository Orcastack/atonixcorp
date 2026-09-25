package hypervisor

import (
	"fmt"
)

type LibvirtDriver struct{}

func NewLibvirtDriver() *LibvirtDriver {
	return &LibvirtDriver{}
}

func (l *LibvirtDriver) StartVM(id string) error {
	fmt.Println("libvirt: starting VM", id)
	// TODO: use libvirt-go bindings
	return nil
}

func (l *LibvirtDriver) StopVM(id string) error {
	fmt.Println("libvirt: stopping VM", id)
	return nil
}

func (l *LibvirtDriver) RebootVM(id string) error {
	fmt.Println("libvirt: rebooting VM", id)
	return nil
}

func (l *LibvirtDriver) GetState(id string) (VMState, error) {
	fmt.Println("libvirt: querying VM state", id)
	return VMRunning, nil
}
