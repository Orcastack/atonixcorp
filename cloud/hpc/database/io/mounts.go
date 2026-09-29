package io

import "fmt"

// Mounts handles dataset mounting for Slurm jobs and Kubernetes pods.
type Mounts struct{}

func NewMounts() *Mounts {
	return &Mounts{}
}

func (m *Mounts) MountForSlurm(dataset string, jobID string) string {
	mountPath := fmt.Sprintf("/mnt/slurm/%s/%s", jobID, dataset)
	fmt.Printf("Mounting dataset=%s for Slurm job=%s at %s\n",
		dataset, jobID, mountPath)
	return mountPath
}

func (m *Mounts) MountForK8s(dataset string, pod string) string {
	mountPath := fmt.Sprintf("/mnt/k8s/%s/%s", pod, dataset)
	fmt.Printf("Mounting dataset=%s for K8s pod=%s at %s\n",
		dataset, pod, mountPath)
	return mountPath
}

func (m *Mounts) Unmount(path string) {
	fmt.Printf("Unmounting dataset at %s\n", path)
}
