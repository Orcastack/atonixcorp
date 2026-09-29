package model

// Metadata stores descriptive and structural information
// about a dataset.
type Metadata struct {
	DatasetID   string
	Dimensions  []int
	Variables   []string
	Attributes  map[string]string
	Compression string
	Encoding    string
	CreatedAt   string
	UpdatedAt   string
}
