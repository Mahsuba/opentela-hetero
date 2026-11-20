package main

// parse.go - turn nvidia-smi text into GPU entries (schema v1).

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reCC   = regexp.MustCompile(`^[0-9]+\.[0-9]+$`)
	reUUID = regexp.MustCompile(`^GPU-([0-9a-fA-F]{8})`)
	rePCI  = regexp.MustCompile(`^[0-9A-Fa-f]+:[0-9A-Fa-f]+:[0-9A-Fa-f]+\.[0-9A-Fa-f]+$`)
)

// parse is all-or-nothing: one bad line -> query_failed, no half lists.
func parse(text string) GPUResult {
	if strings.Contains(text, "No devices were found") {
		return failed(StatusNoneFound, "nvidia-smi", "nvidia-smi reports no devices")
	}
	gpus := []GPU{}
	for n, line := range strings.Split(text, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		f := strings.Split(line, ",")
		if len(f) != 7 {
			return failed(StatusQueryFailed, "nvidia-smi", fmt.Sprintf("line %d: expected 7 values, got %d", n+1, len(f)))
		}
		for i := range f {
			f[i] = strings.TrimSpace(f[i])
		}
		total, err1 := strconv.Atoi(f[1])
		used, err2 := strconv.Atoi(f[2])
		if err1 != nil || err2 != nil || total < 0 || used < 0 {
			return failed(StatusQueryFailed, "nvidia-smi", fmt.Sprintf("line %d: memory is not a whole number", n+1))
		}
		gpus = append(gpus, GPU{Vendor: "nvidia", Name: f[0], MemoryTotalMiB: total, MemoryUsedMiB: used,
			ComputeCapability: opt(f[3], reCC), DriverVersion: opt(f[4], nil),
			UUIDShort: short(f[5]), PCIBusID: opt(f[6], rePCI)})
	}
	if len(gpus) == 0 {
		return failed(StatusNoneFound, "nvidia-smi", "nvidia-smi listed no GPUs")
	}
	return GPUResult{Status: StatusOK, Tool: "nvidia-smi", GPUs: gpus}
}

// opt: "not available" words (or a wrong shape) -> nil -> JSON null.
func opt(v string, re *regexp.Regexp) *string {
	switch v {
	case "", "0", "N/A", "[N/A]", "Not Supported", "[Not Supported]":
		return nil
	}
	if re != nil && !re.MatchString(v) {
		return nil
	}
	return &v
}

func short(v string) *string {
	m := reUUID.FindStringSubmatch(v)
	if m == nil {
		return nil
	}
	s := "GPU-" + strings.ToLower(m[1])
	return &s
}
