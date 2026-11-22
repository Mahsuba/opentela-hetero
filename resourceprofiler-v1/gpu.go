package main

// gpu.go - ask nvidia-smi about this machine's GPUs right now (live mode).
// Fixture files are read in fixture.go.

import (
	"context"
	"os/exec"
	"time"
)

// GPU = one entry of "gpus" in schema v1. nil pointer -> JSON null.
type GPU struct {
	Vendor            string  `json:"vendor"`
	Name              string  `json:"name"`
	MemoryTotalMiB    int     `json:"memory_total_mib"`
	MemoryUsedMiB     int     `json:"memory_used_mib"`
	ComputeCapability *string `json:"compute_capability"`
	DriverVersion     *string `json:"driver_version"`
	UUIDShort         *string `json:"uuid_short"`
	PCIBusID          *string `json:"pci_bus_id"`
}

const query = "--query-gpu=name,memory.total,memory.used,compute_cap,driver_version,uuid,pci.bus_id"

func collectLive() GPUResult {
	if _, err := exec.LookPath("nvidia-smi"); err != nil {
		if _, err2 := exec.LookPath("rocm-smi"); err2 == nil {
			return failed(StatusQueryFailed, "rocm-smi", "rocm-smi found, but AMD reading is not written yet (no AMD GPU to test on)")
		}
		return failed(StatusToolMissing, "none", "neither nvidia-smi nor rocm-smi is installed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "nvidia-smi", query, "--format=csv,noheader,nounits")
	cmd.WaitDelay = 2 * time.Second // don't wait forever on stuck child processes
	out, err := cmd.CombinedOutput()
	if err != nil {
		return classifyRunError("nvidia-smi", err, string(out), ctx.Err())
	}
	return parse(string(out))
}
