package registry

// DatasetMetadata stores descriptive information about a dataset.
type DatasetMetadata struct {
	ID          string
	Name        string
	Format      string // HDF5, NetCDF, ADIOS2, Zarr
	Owner       string
	Tags        []string
	Description string
	CreatedAt   string
	UpdatedAt   string
}
