package state

type InstanceState struct {
	ID      string `json:"id"`
	Desired string `json:"desired"`
	Actual  string `json:"actual"`
	Node    string `json:"node"`
}
