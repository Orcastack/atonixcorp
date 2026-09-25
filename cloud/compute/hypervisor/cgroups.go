package hypervisor

import (
	"fmt"
	"os"
)

type Cgroup struct {
	Path string
}

func NewCgroup(path string) *Cgroup {
	os.MkdirAll(path, 0755)
	return &Cgroup{Path: path}
}

func (c *Cgroup) SetCPU(limit string) error {
	file := c.Path + "/cpu.max"
	fmt.Println("cgroups: set cpu.max =", limit)
	return os.WriteFile(file, []byte(limit), 0644)
}

func (c *Cgroup) SetMemory(limit string) error {
	file := c.Path + "/memory.max"
	fmt.Println("cgroups: set memory.max =", limit)
	return os.WriteFile(file, []byte(limit), 0644)
}
