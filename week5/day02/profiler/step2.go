package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// The same query we ran by hand (see ../fixtures/README.md).
const query = "name,memory.total,memory.used,compute_cap,driver_version"

func main() {
	var text, source string

	if len(os.Args) > 1 {
		// A file was given: read a saved real capture (fixture).
		data, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println("could not read fixture:", err)
			os.Exit(1)
		}
		text = string(data)
		source = "FIXTURE " + os.Args[1]
	} else {
		// No file given: ask the real GPU right now.
		out, err := exec.Command("nvidia-smi",
			"--query-gpu="+query,
			"--format=csv,noheader,nounits").Output()
		if err != nil {
			fmt.Println("source: LIVE nvidia-smi")
			fmt.Println("no NVIDIA GPU found (", err, ")")
			return
		}
		text = string(out)
		source = "LIVE nvidia-smi"
	}

	fmt.Println("source:", source)

	text = strings.TrimSpace(text)
	if text == "" {
		fmt.Println("no GPUs listed")
		return
	}

	// One line = one GPU.
	for i, line := range strings.Split(text, "\n") {
		parts := strings.Split(line, ",")
		if len(parts) != 5 {
			fmt.Println("unexpected line, skipping:", line)
			continue
		}
		for j := range parts {
			parts[j] = strings.TrimSpace(parts[j])
		}
		fmt.Printf("GPU %d\n", i)
		fmt.Println("  name:               ", parts[0])
		fmt.Println("  memory total (MiB): ", parts[1])
		fmt.Println("  memory used (MiB):  ", parts[2])
		fmt.Println("  compute capability: ", parts[3])
		fmt.Println("  driver version:     ", parts[4])
	}
}
