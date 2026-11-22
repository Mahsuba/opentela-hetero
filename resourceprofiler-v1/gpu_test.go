package main

// gpu_test.go - Day 6. Runs collectLive() against FAKE nvidia-smi
// scripts in a temporary folder. These are simulated failures,
// not real broken machines. Needs Linux (/bin/sh).

import (
	"os"
	"path/filepath"
	"testing"
)

// fakeTool puts a fake nvidia-smi in a new empty folder and makes
// that folder the only place the program looks for tools.
func fakeTool(t *testing.T, script string) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nvidia-smi")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
}

func TestNoToolInstalled(t *testing.T) {
	t.Setenv("PATH", t.TempDir()) // empty folder: no nvidia-smi, no rocm-smi
	r := collectLive()
	if r.Status != StatusToolMissing || r.Tool != "none" {
		t.Errorf("want tool_missing + none, got %s + %s", r.Status, r.Tool)
	}
}

func TestFakeNvidiaSmi(t *testing.T) {
	cases := []struct{ name, script, want string }{
		{"works", "echo '" + good + "'", StatusOK},
		{"no devices", "echo 'No devices were found'; exit 6", StatusNoneFound},
		{"no driver", "echo 'NVIDIA-SMI has failed'; exit 9", StatusQueryFailed},
		{"no permission", "echo 'Insufficient Permissions'; exit 4", StatusQueryFailed},
		{"garbage", "echo 'hello world'", StatusQueryFailed},
	}
	for _, c := range cases {
		fakeTool(t, c.script)
		r := collectLive()
		if r.Status != c.want {
			t.Errorf("%s: want %s, got %s (%s)", c.name, c.want, r.Status, r.Err)
		}
	}
}
