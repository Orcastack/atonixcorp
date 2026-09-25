package agent

import (
	"atonixcorp/cloud/compute/hypervisor"
	"atonixcorp/cloud/compute/rpc"
	"fmt"
)

func (a *Agent) HandleMessage(msg rpc.RPCMessage) {
	switch msg.Type {

	case "vm.start":
		vmID := msg.Payload.(string)
		fmt.Println("Agent: start VM", vmID)
		hv := hypervisor.NewLibvirtDriver()
		hv.StartVM(vmID)

	case "vm.pause":
		vmID := msg.Payload.(string)
		fmt.Println("Agent: pause VM", vmID)
		hv := hypervisor.NewLibvirtDriver()
		// TODO: implement pause via libvirt
		_ = hv.StopVM(vmID) // placeholder

	case "vm.stop":
		vmID := msg.Payload.(string)
		fmt.Println("Agent: stop VM", vmID)
		hv := hypervisor.NewLibvirtDriver()
		hv.StopVM(vmID)

	case "vm.reboot":
		vmID := msg.Payload.(string)
		fmt.Println("Agent: reboot VM", vmID)
		hv := hypervisor.NewLibvirtDriver()
		hv.RebootVM(vmID)

	case "vm.pin":
		pin := msg.Payload.(hypervisor.CPUPinning)
		fmt.Println("Agent: CPU pinning", pin)
		hypervisor.PinCPU(pin.VCPU, pin)

	case "vm.numa":
		numa := msg.Payload.(hypervisor.NUMANode)
		fmt.Println("Agent: NUMA placement", numa)
		hypervisor.PlaceNUMA("vm", numa)

	default:
		fmt.Println("Agent: unknown message:", msg.Type)
	}
}
