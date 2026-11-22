package main

// main.go - resourceprofiler v1: print one capability profile (schema v1).
//   resourceprofiler          -> LIVE: this machine, right now
//   resourceprofiler FILE     -> FIXTURE: a saved capture (see fixture.go)
// A node with no GPU is a normal answer, so the exit code is still 0.
// Exit 2 = we were used wrongly (bad arguments or a bad fixture file).

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

const usage = "usage: resourceprofiler [FIXTURE_FILE]"

// build makes the profile. Kept apart from main() so tests can call it.
func build(args []string) (Profile, error) {
	p := Profile{SchemaVersion: "1"}
	var r GPUResult
	if len(args) == 0 {
		p.CollectedAt = time.Now().UTC().Format(timeLayout)
		p.CPU.Arch = runtime.GOARCH
		r = collectLive()
		p.Source = Source{Mode: "live", Tool: r.Tool}
	} else {
		c, err := readFixture(args[0])
		if err != nil {
			return Profile{}, err
		}
		p.CollectedAt, p.CPU.Arch = c.CollectedAt, c.Arch
		r = parse(c.Body)
		p.Source = Source{Mode: "fixture", Tool: r.Tool, FixtureFile: args[0]}
	}
	p.GPUStatus, p.GPUError, p.GPUs = r.Status, r.Err, r.GPUs
	return p, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println(usage)
		return
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	p, err := build(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "resourceprofiler:", err)
		os.Exit(2)
	}
	out, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot build JSON:", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
