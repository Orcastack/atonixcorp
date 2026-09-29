# HDF5 Storage Module

This module provides HDF5-based scientific data storage for the AtonixCorp HPC subsystem.  
It is used for:

- HPC solver outputs (CFD, FEM, physics engines)
- scientific datasets (climate, environmental, orbital mechanics)
- robotics simulation outputs
- aerospace simulation outputs

## Structure

- `hdf5_config.yaml` — configuration for HDF5 chunking, compression, and dataset layout
- `hdf5_example.go` — minimal Go example showing how to initialize HDF5 bindings

## Features

- Hierarchical scientific dataset structure
- Chunked I/O for large datasets
- Optional compression
- Scientific metadata support
- Domain-aware configuration

## Domain Usage

### Robotics
Used for simulation outputs and sensor dataset packaging.

### Aerospace
Used for CFD/FEM solver outputs and mission simulation datasets.

### HPC
Used for timestep storage, residual storage, and scientific variable datasets.

### Scientific
Used for climate, environmental, and orbital mechanics datasets.

## Notes

This folder is part of:
`atonixcorp/cloud/hpc/storage/hdf5/`
