package main

// errors.go - Day 5 "fail safely".
// A node may have no GPU, no GPU tool, or a tool that breaks.
// Each case becomes ONE clear status (from schema v1), never a crash
// and never made-up numbers.

import (
	"context"
	"errors"
	"regexp"
	"strings"
)

// The four values allowed by week5/day04/schema.json (gpu_status).
const (
	StatusOK          = "ok"           // tool ran, at least one GPU read
	StatusNoneFound   = "none_found"   // tool ran fine, but saw no GPU
	StatusToolMissing = "tool_missing" // nvidia-smi / rocm-smi not installed
	StatusQueryFailed = "query_failed" // tool exists but could not answer
)

// GPUResult is what the GPU step hands back to main.go.
type GPUResult struct {
	Status string
	Tool   string // "nvidia-smi", "rocm-smi" or "none"
	GPUs   []GPU
	Err    string // plain reason; empty when Status is "ok"
}

// failed builds a non-ok result. GPUs is always an empty list, so a
// failure can never carry fake or half-read numbers.
func failed(status, tool, why string) GPUResult {
	return GPUResult{Status: status, Tool: tool, GPUs: []GPU{}, Err: safe(why)}
}

// classifyRunError: nvidia-smi exists but exited with an error.
func classifyRunError(tool string, err error, output string, ctxErr error) GPUResult {
	if errors.Is(ctxErr, context.DeadlineExceeded) {
		return failed(StatusQueryFailed, tool, tool+" did not answer within the time limit")
	}
	if strings.Contains(output, "No devices were found") {
		return failed(StatusNoneFound, tool, tool+" reports no devices")
	}
	msg := firstLine(output)
	if msg == "" {
		msg = err.Error()
	}
	return failed(StatusQueryFailed, tool, tool+" failed: "+msg)
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

// safe cuts any full GPU uuid down to "GPU-" + 8 hex, so an error
// message can never leak the full id into the JSON (privacy rule W5D3).
var fullUUID = regexp.MustCompile(`(GPU-[0-9a-fA-F]{8})[-0-9a-fA-F]+`)

func safe(s string) string { return fullUUID.ReplaceAllString(s, "$1") }
