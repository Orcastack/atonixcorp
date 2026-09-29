package registry

import "fmt"

// DatasetRegistry is the central catalog for all datasets.
type DatasetRegistry struct {
	datasets map[string]*DatasetMetadata
}

func NewDatasetRegistry() *DatasetRegistry {
	return &DatasetRegistry{
		datasets: make(map[string]*DatasetMetadata),
	}
}

func (r *DatasetRegistry) Register(meta *DatasetMetadata) {
	r.datasets[meta.ID] = meta
	fmt.Printf("Dataset registered: %s\n", meta.ID)
}

func (r *DatasetRegistry) Get(id string) *DatasetMetadata {
	return r.datasets[id]
}

func (r *DatasetRegistry) List() []*DatasetMetadata {
	list := []*DatasetMetadata{}
	for _, d := range r.datasets {
		list = append(list, d)
	}
	return list
}

func (r *DatasetRegistry) Delete(id string) {
	delete(r.datasets, id)
	fmt.Printf("Dataset deleted: %s\n", id)
}
