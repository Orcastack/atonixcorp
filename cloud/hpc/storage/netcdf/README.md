# NetCDF Storage Module

This module provides NetCDF-based scientific data storage for the AtonixCorp HPC subsystem.  
NetCDF is widely used for:

- climate and environmental datasets
- atmospheric and oceanographic models
- aerospace mission environmental data
- robotics environmental mapping
- HPC scientific simulations

## Structure

- `netcdf_config.yaml` — configuration for dimensions, variables, compression, and metadata
- `netcdf_example.go` — minimal Go example showing how to initialize NetCDF bindings

## Features

- Multi-dimensional scientific datasets
- Strong metadata support (units, dimensions, attributes)
- Domain-aware variable definitions
- Optional compression
- High compatibility with scientific tools

## Domain Usage

### Robotics
Used for environmental mapping, sensor fusion datasets, and simulation outputs.

### Aerospace
Used for mission environmental data, atmospheric models, and CFD boundary conditions.

### HPC
Used for climate models, physics simulations, and multi-dimensional solver outputs.

### Scientific
Used for Earth science, weather models, and environmental datasets.

## Notes

This folder is part of:
`atonixcorp/cloud/hpc/storage/netcdf/`
