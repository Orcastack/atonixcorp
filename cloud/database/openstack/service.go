package openstack

type ServiceClient struct {
	compute *ComputeClient
	volume  *VolumeClient
}

func NewServiceClient(client *Client) *ServiceClient {
	return &ServiceClient{
		compute: NewComputeClient(client),
		volume:  NewVolumeClient(client),
	}
}

func (c *ServiceClient) CreateServer(name string, cpu int, memoryMB int, networkID string, securityGroupID string) (string, error) {
	return c.compute.CreateServer(name, cpu, memoryMB, networkID, securityGroupID)
}

func (c *ServiceClient) CreateVolume(name string, sizeGB int) (string, error) {
	return c.volume.CreateVolume(name, sizeGB)
}

func (c *ServiceClient) AttachVolume(serverID, volumeID string) error {
	return c.volume.AttachVolume(serverID, volumeID)
}

func (c *ServiceClient) GetServerIP(serverID string) (string, error) {
	return c.compute.GetServerIP(serverID)
}
