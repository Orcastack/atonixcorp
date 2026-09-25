package hypervisor

import (
	"fmt"
	"os/exec"
)

func QemuStart(binary string, args ...string) error {
	fmt.Println("qemu: exec", binary, args)
	cmd := exec.Command(binary, args...)
	return cmd.Start()
}

func QemuStop(pid int) error {
	fmt.Println("qemu: kill pid", pid)
	return exec.Command("kill", fmt.Sprintf("%d", pid)).Run()
}
