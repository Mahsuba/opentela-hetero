package main

// parse_test.go - Day 6. Checks parse() with saved text, so no GPU is needed.
// Inputs marked "made-up" are test text, not real hardware readings.

import "testing"

// val prints a *string for error messages (nil -> "null").
func val(p *string) string {
	if p == nil {
		return "null"
	}
	return *p
}

// One made-up good line, reused below.
const good = "Test GPU, 100, 10, 7.5, 1.0, GPU-12345678, 00000000:01:00.0"

func TestParseCases(t *testing.T) {
	cases := []struct {
		name, input, want string
		gpus              int
	}{
		{"one good line", good, StatusOK, 1},
		{"blank lines are skipped", "\n" + good + "\n\n", StatusOK, 1},
		{"empty output", "", StatusNoneFound, 0},
		{"tool says no devices", "No devices were found", StatusNoneFound, 0},
		{"only 6 values", "Test GPU, 100, 10, 7.5, 1.0, GPU-12345678", StatusQueryFailed, 0},
		{"memory is text", "Test GPU, lots, 10, 7.5, 1.0, GPU-12345678, 00000000:01:00.0", StatusQueryFailed, 0},
		{"negative memory", "Test GPU, -5, 10, 7.5, 1.0, GPU-12345678, 00000000:01:00.0", StatusQueryFailed, 0},
		{"used more than total", "Test GPU, 100, 200, 7.5, 1.0, GPU-12345678, 00000000:01:00.0", StatusQueryFailed, 0},
		{"good line then bad line", good + "\nhello world", StatusQueryFailed, 0},
		{"known limit: comma in name", "Test, GPU, 100, 10, 7.5, 1.0, GPU-12345678, 00000000:01:00.0", StatusQueryFailed, 0},
	}
	for _, c := range cases {
		r := parse(c.input)
		if r.Status != c.want || len(r.GPUs) != c.gpus {
			t.Errorf("%s: want %s + %d GPUs, got %s + %d (%s)", c.name, c.want, c.gpus, r.Status, len(r.GPUs), r.Err)
		}
	}
}

// [N/A] must become null, and a full uuid must be cut to GPU- + 8 hex.
func TestNotAvailableAndPrivacy(t *testing.T) {
	r := parse("Test GPU, 100, 0, [N/A], [N/A], GPU-ABCDEF12-1111-2222-3333-444455556666, 00000000:00:04.0")
	if r.Status != StatusOK || len(r.GPUs) != 1 {
		t.Fatalf("want ok + 1 GPU, got %s (%s)", r.Status, r.Err)
	}
	g := r.GPUs[0]
	if g.ComputeCapability != nil || g.DriverVersion != nil {
		t.Errorf("[N/A] should be null, got %s and %s", val(g.ComputeCapability), val(g.DriverVersion))
	}
	if val(g.UUIDShort) != "GPU-abcdef12" {
		t.Errorf("uuid_short: want GPU-abcdef12, got %s", val(g.UUIDShort))
	}
}
