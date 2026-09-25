package hypervisor

import "fmt"

type CPUPinning struct {
	VCPU  int
	PCPUs []int
}

func PinCPU(vmID string, pin CPUPinning) error {
	fmt.Printf("cpu: pinning vCPU %d to pCPUs %v for VM %s\n",
		pin.VCPU, pin.PCPUs, vmID)
	// TODO: write to cpuset cgroup
	return nil
}
