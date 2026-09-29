package model

// Lineage tracks dataset origin, transformations, and dependencies.
type Lineage struct {
	DatasetID      string
	Parents        []string
	Children       []string
	TransformSteps []string
	CreatedAt      string
	UpdatedAt      string
}
