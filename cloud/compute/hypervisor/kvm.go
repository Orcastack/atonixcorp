package hypervisor

import "fmt"

type KVM struct{}

func NewKVM() *KVM {
	return &KVM{}
}

func (k *KVM) StartVM(id string) error {
	fmt.Println("KVM: start VM", id)
	return nil
}

func (k *KVM) StopVM(id string) error {
	fmt.Println("KVM: stop VM", id)
	return nil
}

func (k *KVM) RebootVM(id string) error {
	fmt.Println("KVM: reboot VM", id)
	return nil
}

func (k *KVM) GetState(id string) (VMState, error) {
	fmt.Println("KVM: get state for VM", id)
	return VMRunning, nil
}
