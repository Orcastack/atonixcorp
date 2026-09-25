package datacenter

import (
	"fmt"

	"atnetwork/core"
)

type FabricBackend struct {
	vlan  *VLANManager
	vxlan *VXLANManager
}

func NewFabricBackend() *FabricBackend {
	return &FabricBackend{
		vlan:  NewVLANManager(),
		vxlan: NewVXLANManager(),
	}
}

func (b *FabricBackend) ApplyNetwork(n core.Network) error {
	fmt.Println("DC fabric: ApplyNetwork", n.Name, n.CIDR)

	// Example: choose VLAN or VXLAN based on config
	if err := b.vlan.AllocateForNetwork(n); err != nil {
		return err
	}

	return nil
}

func (b *FabricBackend) ApplySubnet(s core.Subnet) error {
	fmt.Println("DC fabric: ApplySubnet", s.CIDR)
	return nil
}

func (b *FabricBackend) ApplyRouter(r core.Router) error {
	fmt.Println("DC fabric: ApplyRouter", r.Name)
	return nil
}

func (b *FabricBackend) ApplyPort(p core.Port) error {
	fmt.Println("DC fabric: ApplyPort", p.ID)
	return nil
}

func (b *FabricBackend) ApplySecurityGroup(sg core.SecurityGroup) error {
	fmt.Println("DC fabric: ApplySecurityGroup", sg.Name)
	return nil
}
