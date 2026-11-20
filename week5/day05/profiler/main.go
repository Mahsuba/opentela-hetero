package main

// main.go - print one capability profile (schema v1) as JSON.
//   no argument  -> LIVE (asks nvidia-smi now)
//   file argument -> FIXTURE (reads saved nvidia-smi output)
// A node with no GPU is a normal answer, so exit code is 0.

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"
)

type Source struct {
	Mode        string `json:"mode"`
	Tool        string `json:"tool"`
	FixtureFile string `json:"fixture_file,omitempty"`
}

type CPU struct {
	Arch string `json:"arch"`
}

type Profile struct {
	SchemaVersion string `json:"schema_version"`
	CollectedAt   string `json:"collected_at"`
	Source        Source `json:"source"`
	CPU           CPU    `json:"cpu"`
	GPUStatus     string `json:"gpu_status"`
	GPUError      string `json:"gpu_error,omitempty"`
	GPUs          []GPU  `json:"gpus"`
}

func main() {
	p := Profile{
		SchemaVersion: "1",
		CollectedAt:   time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		CPU:           CPU{Arch: runtime.GOARCH},
	}
	var r GPUResult
	if len(os.Args) > 1 {
		r = collectFixture(os.Args[1])
		p.Source = Source{Mode: "fixture", Tool: r.Tool, FixtureFile: os.Args[1]}
	} else {
		r = collectLive()
		p.Source = Source{Mode: "live", Tool: r.Tool}
	}
	p.GPUStatus, p.GPUError, p.GPUs = r.Status, r.Err, r.GPUs

	out, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot build JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
