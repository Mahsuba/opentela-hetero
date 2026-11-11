package main

import "fmt"

type Fabric struct {
	Nodes      map[string]Node
	KnownNodes map[string]map[string]Node
}

func (fabric *Fabric) exchangeCapability(message Message) {
	sender := fabric.Nodes[message.From]

	receiverKnowledge := fabric.KnownNodes[message.To]

	if receiverKnowledge == nil {
		receiverKnowledge = make(map[string]Node)
		fabric.KnownNodes[message.To] = receiverKnowledge
	}

	receiverKnowledge[sender.ID] = sender

	fmt.Printf(
		"%s received capability information from %s\n",
		message.To,
		message.From,
	)
}

func (fabric *Fabric) showKnowledge() {
	fmt.Println("\n=== Node Knowledge ===")

	for nodeID, knownNodes := range fabric.KnownNodes {
		fmt.Printf("\n%s knows about:\n", nodeID)

		for knownID, node := range knownNodes {
			fmt.Printf(
				"  %s | %s | %s | %d GB\n",
				knownID,
				node.Location,
				node.GPU,
				node.VRAM,
			)
		}
	}
}

func (fabric *Fabric) detectFailure(nodeID string) {
	node := fabric.Nodes[nodeID]

	node.Status = "suspected"

	fabric.Nodes[nodeID] = node

	fmt.Printf("%s is suspected to have failed\n", nodeID)
}
