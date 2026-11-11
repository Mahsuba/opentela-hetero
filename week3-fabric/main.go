package main

import "fmt"

func main() {
	fabric := Fabric{
		Nodes: map[string]Node{
			"node-a": {
				ID:       "node-a",
				Location: "Dhaka",
				GPU:      "RTX 5060",
				VRAM:     8,
				Status:   "online",
			},
			"node-b": {
				ID:       "node-b",
				Location: "London",
				GPU:      "RTX 5070",
				VRAM:     12,
				Status:   "online",
			},
			"node-c": {
				ID:       "node-c",
				Location: "Zurich",
				GPU:      "RTX 5090",
				VRAM:     32,
				Status:   "online",
			},
		},

		KnownNodes: make(map[string]map[string]Node),
	}

	fmt.Println("\n=== Capability Exchange ===")

	fabric.exchangeCapability(Message{
		From: "node-a",
		To:   "node-b",
	})

	fabric.exchangeCapability(Message{
		From: "node-b",
		To:   "node-c",
	})

	fabric.exchangeCapability(Message{
		From: "node-c",
		To:   "node-a",
	})

	fabric.showKnowledge()

	fmt.Println("\n=== Failure Detection ===")

	fabric.detectFailure("node-b")

	if fabric.Nodes["node-b"].Status == "suspected" {
		fmt.Println("SUCCESS: node-b failure detected correctly")
	} else {
		fmt.Println("ERROR: node-b status was not updated")
	}
	fmt.Println("=== Tiny Fabric ===")

	for id, node := range fabric.Nodes {
		fmt.Printf(
			"%s | %s | %s | %d GB | %s\n",
			id,
			node.Location,
			node.GPU,
			node.VRAM,
			node.Status,
		)
	}
}
