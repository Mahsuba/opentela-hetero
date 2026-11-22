package main

// fixture_test.go - Day 7. Fixture files and the whole profile (build).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFixture saves made-up text as a fixture file in a temp folder.
func writeFixture(t *testing.T, text string) string {
	path := filepath.Join(t.TempDir(), "f.csv")
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// Real capture from Kaggle "GPU T4 x2" (see testdata/README.md).
func TestKaggleFixture(t *testing.T) {
	p, err := build([]string{"testdata/kaggle-t4x2.csv"})
	if err != nil {
		t.Fatal(err)
	}
	// The Day 5 bug: these must come from the file, not from this machine.
	if p.CollectedAt != "2026-11-20T20:37:34Z" || p.CPU.Arch != "amd64" {
		t.Errorf("want capture time + amd64, got %s + %s", p.CollectedAt, p.CPU.Arch)
	}
	if p.GPUStatus != StatusOK || len(p.GPUs) != 2 || p.Source.Mode != "fixture" {
		t.Fatalf("want fixture ok + 2 GPUs, got %s %s + %d", p.Source.Mode, p.GPUStatus, len(p.GPUs))
	}
	for _, g := range p.GPUs {
		if g.Name != "Tesla T4" || g.MemoryTotalMiB != 15360 || val(g.ComputeCapability) != "7.5" {
			t.Errorf("wrong values: %+v", g)
		}
	}
	a, b := p.GPUs[0], p.GPUs[1]
	if val(a.UUIDShort) != "GPU-b6f3fe4c" || val(b.UUIDShort) != "GPU-57dc86ae" {
		t.Errorf("uuid_short: got %s and %s", val(a.UUIDShort), val(b.UUIDShort))
	}
	if val(a.PCIBusID) == val(b.PCIBusID) {
		t.Errorf("two cards must have different pci_bus_id")
	}
}

// A made-up arm64 capture must say arm64 even when replayed on an amd64 PC.
func TestFixtureKeepsItsOwnArch(t *testing.T) {
	p, err := build([]string{writeFixture(t, "# collected_at=2026-01-02T03:04:05Z\n# arch=aarch64\n")})
	if err != nil {
		t.Fatal(err)
	}
	if p.CPU.Arch != "arm64" || p.CollectedAt != "2026-01-02T03:04:05Z" {
		t.Errorf("want arm64 + file time, got %s + %s", p.CPU.Arch, p.CollectedAt)
	}
	if p.GPUStatus != StatusNoneFound || len(p.GPUs) != 0 {
		t.Errorf("no GPU lines: want none_found, got %s", p.GPUStatus)
	}
}

// Bad fixture files are refused (error, no profile), never guessed.
func TestBadFixtureRefused(t *testing.T) {
	gpu := "\nTest GPU, 100, 10, 7.5, 1.0, GPU-12345678, 00000000:01:00.0\n"
	cases := []struct{ name, text, reason string }{
		{"no header at all", gpu, "collected_at"},
		{"no arch", "# collected_at=2026-01-02T03:04:05Z" + gpu, "arch"},
		{"no time", "# arch=amd64" + gpu, "collected_at"},
		{"time not UTC format", "# collected_at=2026-01-02 03:04\n# arch=amd64" + gpu, "collected_at"},
		{"unknown arch", "# collected_at=2026-01-02T03:04:05Z\n# arch=mips" + gpu, "arch"},
	}
	for _, c := range cases {
		_, err := build([]string{writeFixture(t, c.text)})
		if err == nil || !strings.Contains(err.Error(), c.reason) {
			t.Errorf("%s: want error about %s, got %v", c.name, c.reason, err)
		}
	}
}

func TestMissingFixtureFile(t *testing.T) {
	_, err := build([]string{"testdata/does-not-exist.csv"})
	if err == nil || !strings.Contains(err.Error(), "cannot read fixture") {
		t.Errorf("want cannot read fixture, got %v", err)
	}
}

// Error line numbers must match the file, "#" lines included.
func TestErrorLineNumber(t *testing.T) {
	p, err := build([]string{writeFixture(t, "# collected_at=2026-01-02T03:04:05Z\n# arch=amd64\nhello world\n")})
	if err != nil {
		t.Fatal(err)
	}
	if p.GPUStatus != StatusQueryFailed || !strings.Contains(p.GPUError, "line 3") {
		t.Errorf("want query_failed on line 3, got %s (%s)", p.GPUStatus, p.GPUError)
	}
}
