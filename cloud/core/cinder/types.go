package cinder

type Volume struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"` // GB
	Status      string `json:"status"`
	Description string `json:"description"`
}
