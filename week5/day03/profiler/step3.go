// step3.go - Week 5 Day 3: device identity
//
// Day 2 asked "WHAT model is this GPU?"  (name, memory)
// Day 3 asks  "WHICH physical card is this?" (uuid, PCI address)
//
// Usage: "go run step3.go" = LIVE (asks nvidia-smi on this machine)
// and "go run step3.go <file.csv>" = FIXTURE (reads a saved capture)
// File line format (no header), same order as the query below:
// index, pci.bus_id, uuid, serial, name, memory.used
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const query = "index,pci.bus_id,uuid,serial,name,memory.used"

func main() {
	var raw []byte
	var err error
	var source string

	if len(os.Args) > 1 {
		source = "FIXTURE " + os.Args[1]
		raw, err = os.ReadFile(os.Args[1])
	} else {
		source = "LIVE nvidia-smi"
		raw, err = exec.Command("nvidia-smi",
			"--query-gpu="+query,
			"--format=csv,noheader,nounits").Output()
	}
	fmt.Println("source:", source)
	if err != nil {
		fmt.Printf("no NVIDIA GPU identity found (%v)\n", err)
		return
	}

	for _, line := range strings.Split(strings.TrimSpace(string(raw)), "\n") {
		f := strings.Split(line, ",")
		if len(f) != 6 {
			fmt.Printf("skipping line: expected 6 fields, got %d: %q\n", len(f), line)
			continue
		}
		for i := range f {
			f[i] = strings.TrimSpace(f[i])
		}
		index, pci, uuid, serial, name, used := f[0], f[1], f[2], f[3], f[4], f[5]

		fmt.Println()
		fmt.Println("GPU", index, "-", name, "(index = order number only)")
		fmt.Println("  IDENTITY (should be the same every run)")
		fmt.Println("    uuid:  ", show(short(clean(uuid)), uuid))
		fmt.Println("    pci:   ", show(clean(pci), pci))
		fmt.Println("    serial:", show(clean(serial), serial))
		fmt.Println("    key:   ", key(clean(uuid), clean(pci), index))
		fmt.Println("  USAGE (expected to change between runs)")
		fmt.Println("    used memory:", used, "MiB")
	}
}

// clean returns "" when the tool gave a placeholder instead of a real value.
// Our RTX 5060 reported serial "0" - no real card has serial number 0.
func clean(s string) string {
	switch strings.ToUpper(s) {
	case "", "0", "N/A", "[N/A]", "NOT SUPPORTED", "[NOT SUPPORTED]":
		return ""
	}
	return s
}

// show prints the value, or says honestly that it is not available.
func show(value, raw string) string {
	if value == "" {
		return fmt.Sprintf("not available (tool said %q)", raw)
	}
	return value
}

// short keeps only the start of the uuid, e.g. "GPU-9877209e".
// Enough to tell cards apart; the full uuid stays on this machine.
func short(u string) string {
	if len(u) > 12 {
		return u[:12] + "..."
	}
	return u
}

// key picks the most stable identifier that is available:
// uuid (stays with the card) > PCI address (stays with the slot) > index (order only).
func key(uuid, pci, index string) string {
	if uuid != "" {
		return "uuid " + short(uuid) + " (best: stays with the card)"
	}
	if pci != "" {
		return "pci " + pci + " (ok: changes if the card moves slot)"
	}
	return "index " + index + " (weak: order number only)"
}
