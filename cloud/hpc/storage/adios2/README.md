# ADIOS2 Storage Module

This module provides ADIOS2-based high-performance I/O integration for the AtonixCorp HPC subsystem.  
It is used for:

- HPC solver streaming (CFD, FEM, physics engines)
- robotics sensor streaming (LiDAR, radar, camera)
- aerospace telemetry streaming
- scientific variable streaming
- distributed parallel I/O

## Structure

- `adios2_config.yaml` — configuration for ADIOS2 engines and transports
- `adios2_example.go` — minimal Go example showing how to initialize ADIOS2 bindings
- `adios2.go` (optional) — module initializer for HPC integration

## Supported Engines

- **SST** — low-latency streaming
- **BP5** — high-throughput batch output
- **DataMan** — distributed streaming

## Domain Usage

### Robotics
Used for real-time sensor ingestion and fusion.

### Aerospace
Used for mission telemetry and solver outputs.

### HPC
Used for timestep streaming, residual streaming, and parallel solver output.

### AI/ML
Used for streaming training batches from HPC simulations.

## Notes

This folder is part of:
`atonixcorp/cloud/hpc/storage/adios2/`
