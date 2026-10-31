package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Resource Profiler v0")
	fmt.Fprintln(os.Stderr, "Checking nodes...")

	nodes := []Node{
		{
			Name:     "computer-1",
			Location: "Dhaka",
			GPU:      "RTX 5060",
			VRAM:     8,
		},
		{
			Name:     "computer-2",
			Location: "London",
			GPU:      "RTX 5070",
			VRAM:     12,
		},
		{
			Name:     "computer-3",
			Location: "Zurich",
			GPU:      "RTX 5090",
			VRAM:     32,
		},
	}

	data, err := json.MarshalIndent(nodes, "", "  ")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	hasErrors := false

	for _, node := range nodes {
		err := validateNode(node)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Validation error in", node.Name+":", err)
			hasErrors = true
		}
	}

	if hasErrors {
		return
	}

	fmt.Println(string(data))
}
