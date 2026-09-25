package neutron

type Network struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	AdminStateUp bool   `json:"admin_state_up"`
	Status       string `json:"status"`
}

type Subnet struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	NetworkID string `json:"network_id"`
	CIDR      string `json:"cidr"`
	IPVersion int    `json:"ip_version"`
}

type Port struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	NetworkID string `json:"network_id"`
	DeviceID  string `json:"device_id"`
}
