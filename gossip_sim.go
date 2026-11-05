package main

import "fmt"

type Node struct {
	Name      string
	Location  string
	KnowsFact bool
	Peers     []int
}

func main() {
	nodes := []Node{
		{
			Name:      "Node A",
			Location:  "Dhaka",
			KnowsFact: true,
			Peers:     []int{1, 2},
		},
		{
			Name:      "Node B",
			Location:  "London",
			KnowsFact: false,
			Peers:     []int{0, 3},
		},
		{
			Name:      "Node C",
			Location:  "Zurich",
			KnowsFact: false,
			Peers:     []int{0, 3},
		},
		{
			Name:      "Node D",
			Location:  "Singapore",
			KnowsFact: false,
			Peers:     []int{1, 2},
		},
	}

	fmt.Println("Gossip Protocol Simulator")
	fmt.Println("--------------------------")
	fmt.Println("Total nodes:", len(nodes))
	// Simulate gossip for 3 rounds.
	for round := 1; round <= 3; round++ {

		fmt.Println("\nRound", round)

		// Remember who knew the fact at the start.
		knownAtStart := make([]bool, len(nodes))

		for i := range nodes {
			knownAtStart[i] = nodes[i].KnowsFact
		}

		// Track newly informed nodes.
		newlyLearned := make([]bool, len(nodes))

		// Spread information to peers.
		for i, node := range nodes {

			if !knownAtStart[i] {
				continue
			}

			for _, peer := range node.Peers {

				if !knownAtStart[peer] && !newlyLearned[peer] {

					newlyLearned[peer] = true

					fmt.Println(
						node.Name,
						"shared the fact with",
						nodes[peer].Name,
					)
				}
			}
		}

		// Apply new knowledge after the round finishes.
		for i := range nodes {
			if newlyLearned[i] {
				nodes[i].KnowsFact = true
			}
		}

		// Display knowledge after this round.
		for _, node := range nodes {
			fmt.Println(node.Name, "knows fact:", node.KnowsFact)
		}
	}
	fmt.Println("\nGossip Spread Diagram")
	fmt.Println("---------------------")
	fmt.Println("Initial: A")
	fmt.Println("Round 1: A -> B, C")
	fmt.Println("Round 2: B -> D (first delivery)")
	fmt.Println("Final:   A, B, C, D")
}
