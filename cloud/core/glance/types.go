package glance

type Image struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	DiskFormat      string `json:"disk_format"`
	ContainerFormat string `json:"container_format"`
	Size            int64  `json:"size"`
	Visibility      string `json:"visibility"`
}
