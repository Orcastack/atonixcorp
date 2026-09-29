package state

type NetworkState struct {
	ID      string `json:"id"`
	Desired string `json:"desired"`
	Actual  string `json:"actual"`
	CIDR    string `json:"cidr"`
}
