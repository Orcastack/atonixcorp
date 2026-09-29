package model

// Dataset is the canonical representation of a scientific dataset
// inside the AtonixCorp HPC Database.
type Dataset struct {
	ID          string
	Name        string
	Format      string // HDF5, NetCDF, ADIOS2, Zarr
	Path        string // physical or logical path
	Owner       string
	Tags        []string
	Description string
	CreatedAt   string
	UpdatedAt   string
}
