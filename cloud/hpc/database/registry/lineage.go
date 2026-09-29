package registry

import "fmt"

// Lineage tracks dataset origin, transformations, and dependencies.
type Lineage struct {
	DatasetID       string
	Parents         []string
	Transformations []string
	Children        []string
}

func NewLineage(id string) *Lineage {
	return &Lineage{
		DatasetID: id,
	}
}

func (l *Lineage) AddParent(parent string) {
	l.Parents = append(l.Parents, parent)
	fmt.Printf("Lineage: %s <- %s\n", l.DatasetID, parent)
}

func (l *Lineage) AddChild(child string) {
	l.Children = append(l.Children, child)
	fmt.Printf("Lineage: %s -> %s\n", l.DatasetID, child)
}

func (l *Lineage) AddTransformation(step string) {
	l.Transformations = append(l.Transformations, step)
}
