package main

// fixture.go - read a saved capture file (fixture mode).
// A fixture says WHEN and on WHICH CPU it was captured, in "#" lines
// at the top, followed by the raw nvidia-smi lines:
//
//	# collected_at=2026-11-20T20:37:34Z
//	# arch=x86_64
//	Tesla T4, 15360, 0, 7.5, 580.178.04, GPU-b6f3fe4c, 00000000:00:04.0
//
// Without these two lines we refuse the file: we never use the time or
// CPU of the machine that happens to replay it (bug found on Day 5).

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const timeLayout = "2006-01-02T15:04:05Z"

// Go names allowed by the schema, plus the two names uname -m prints.
var archNames = map[string]string{
	"amd64": "amd64", "x86_64": "amd64",
	"arm64": "arm64", "aarch64": "arm64",
	"riscv64": "riscv64", "386": "386", "arm": "arm",
	"ppc64le": "ppc64le", "s390x": "s390x",
}

// Capture = what one fixture file says.
type Capture struct {
	CollectedAt string // capture time, UTC
	Arch        string // Go name, e.g. amd64
	Body        string // nvidia-smi lines; each "#" line is left blank
}

func readFixture(path string) (Capture, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Capture{}, fmt.Errorf("cannot read fixture: %v", err)
	}
	var c Capture
	var body []string
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(line, "#") {
			body = append(body, line)
			continue
		}
		body = append(body, "") // keep line numbers in errors = file line numbers
		key, value, _ := strings.Cut(strings.TrimSpace(line[1:]), "=")
		switch strings.TrimSpace(key) {
		case "collected_at":
			c.CollectedAt = strings.TrimSpace(value)
		case "arch":
			c.Arch = strings.TrimSpace(value)
		} // any other "#" line is a free-text note
	}
	if _, err := time.Parse(timeLayout, c.CollectedAt); err != nil {
		return Capture{}, fmt.Errorf("fixture needs a line like \"# collected_at=2026-11-20T20:37:34Z\" (got %q)", c.CollectedAt)
	}
	arch, ok := archNames[c.Arch]
	if !ok {
		return Capture{}, fmt.Errorf("fixture needs a line like \"# arch=x86_64\" (got %q)", c.Arch)
	}
	c.Arch = arch
	c.Body = strings.Join(body, "\n")
	return c, nil
}
