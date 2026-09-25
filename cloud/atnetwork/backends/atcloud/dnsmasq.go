package atcloud

import (
	"atnetwork/core"
	"fmt"
	"os"
)

type DnsmasqBackend struct {
	ConfigDir string
}

func NewDnsmasqBackend(dir string) *DnsmasqBackend {
	return &DnsmasqBackend{ConfigDir: dir}
}

func (b *DnsmasqBackend) ApplySubnet(s core.Subnet) error {
	fmt.Println("ATN → dnsmasq: configure DHCP for", s.CIDR)

	cfg := fmt.Sprintf("dhcp-range=%s,static", s.CIDR)
	path := fmt.Sprintf("%s/%s.conf", b.ConfigDir, s.ID)

	return os.WriteFile(path, []byte(cfg), 0644)
}
