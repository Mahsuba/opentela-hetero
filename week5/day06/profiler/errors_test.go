package main

// errors_test.go - Day 6. Checks how a broken nvidia-smi is labelled.
// All error texts here are made-up test text.

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestClassifyRunError(t *testing.T) {
	exit := errors.New("exit status 9")
	cases := []struct {
		name, output string
		ctxErr       error
		want, reason string
	}{
		{"took too long", "", context.DeadlineExceeded, StatusQueryFailed, "time limit"},
		{"no devices", "No devices were found", nil, StatusNoneFound, "no devices"},
		{"no driver", "NVIDIA-SMI has failed\nmore text", nil, StatusQueryFailed, "NVIDIA-SMI has failed"},
		{"no output at all", "", nil, StatusQueryFailed, "exit status 9"},
	}
	for _, c := range cases {
		r := classifyRunError("nvidia-smi", exit, c.output, c.ctxErr)
		if r.Status != c.want || !strings.Contains(r.Err, c.reason) || len(r.GPUs) != 0 {
			t.Errorf("%s: want %s with %q, got %s with %q", c.name, c.want, c.reason, r.Status, r.Err)
		}
		if strings.Contains(r.Err, "more text") {
			t.Errorf("%s: only the first line should be kept", c.name)
		}
	}
}

// An error message must never carry a full GPU uuid (privacy rule W5D3).
func TestSafeCutsFullUUID(t *testing.T) {
	r := classifyRunError("nvidia-smi", errors.New("exit status 15"),
		"Unknown Error on GPU-9877209e-1111-2222-3333-444455556666", nil)
	if strings.Contains(r.Err, "1111") || !strings.Contains(r.Err, "GPU-9877209e") {
		t.Errorf("uuid not cut correctly: %q", r.Err)
	}
}

// A failure must print "gpus": [] (schema wants a list), never "gpus": null.
func TestFailureGivesEmptyList(t *testing.T) {
	out, err := json.Marshal(failed(StatusToolMissing, "none", "test"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"GPUs":[]`) {
		t.Errorf("want an empty list, got %s", out)
	}
}
