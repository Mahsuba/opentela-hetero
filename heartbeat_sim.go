package main

import (
	"fmt"
	"time"
)

// Node represents one computer in our distributed system.
type Node struct {
	Name                 string
	City                 string
	IntentionallyStopped bool
	HeartbeatConfirmed   bool
	LastHeartbeat        time.Time
}

// A node is suspected unavailable if no heartbeat arrives
// within this time limit.
const heartbeatTimeout = 3 * time.Second

func main() {
	fmt.Println("======================================")
	fmt.Println("       DISTRIBUTED SYSTEM")
	fmt.Println("        HEARTBEAT MONITOR")
	fmt.Println("======================================")

	// Create three computers in different locations.
	// The third computer is intentionally stopped.
	nodes := []Node{
		{
			Name:                 "computer-1",
			City:                 "Dhaka",
			IntentionallyStopped: false,
			HeartbeatConfirmed:   true,
			LastHeartbeat:        time.Now(),
		},
		{
			Name:                 "computer-2",
			City:                 "London",
			IntentionallyStopped: false,
			HeartbeatConfirmed:   true,
			LastHeartbeat:        time.Now(),
		},
		{
			Name:                 "computer-3",
			City:                 "Zurich",
			IntentionallyStopped: true,
			HeartbeatConfirmed:   true,
			LastHeartbeat:        time.Now(),
		},
	}

	// Display the initial configuration.
	fmt.Println("\nINITIAL NODE CONFIGURATION")
	fmt.Println("--------------------------------------")

	for _, node := range nodes {
		fmt.Printf("Computer: %s | Location: %s\n",
			node.Name, node.City)

		if node.IntentionallyStopped {
			fmt.Println("Configured behavior: Intentionally stopped")
		} else {
			fmt.Println("Configured behavior: Running normally")
		}

		fmt.Println("Initial heartbeat: Previously confirmed")
		fmt.Println("--------------------------------------")
	}

	// Start monitoring heartbeats.
	fmt.Println("\nSTARTING HEARTBEAT MONITORING")
	fmt.Printf("Heartbeat interval: 1 second\n")
	fmt.Printf("Heartbeat timeout: %v\n", heartbeatTimeout)

	for round := 1; round <= 6; round++ {
		fmt.Printf("\n========== MONITORING ROUND %d ==========\n", round)

		for i := range nodes {

			// Running computers send a new heartbeat.
			if !nodes[i].IntentionallyStopped {
				nodes[i].LastHeartbeat = time.Now()
			}

			// Calculate the time since the last heartbeat.
			elapsed := time.Since(nodes[i].LastHeartbeat)

			// Check whether the heartbeat has timed out.
			if elapsed > heartbeatTimeout {
				nodes[i].HeartbeatConfirmed = false

				fmt.Printf(
					"WARNING: %s (%s) is SUSPECTED UNAVAILABLE.\n",
					nodes[i].Name,
					nodes[i].City,
				)

				fmt.Printf(
					"Reason: No heartbeat received for %v.\n",
					elapsed.Round(time.Millisecond),
				)

			} else if !nodes[i].IntentionallyStopped {
				nodes[i].HeartbeatConfirmed = true

				fmt.Printf(
					"HEARTBEAT RECEIVED: %s (%s) is responding normally.\n",
					nodes[i].Name,
					nodes[i].City,
				)

			} else {
				fmt.Printf(
					"WAITING: No new heartbeat from %s (%s).\n",
					nodes[i].Name,
					nodes[i].City,
				)
			}
		}

		time.Sleep(1 * time.Second)
	}

	// Allow additional time for the stopped computer's
	// heartbeat to exceed the timeout threshold.
	fmt.Println("\nWaiting for remaining heartbeat timeouts...")
	time.Sleep(2 * time.Second)

	// Display the final monitoring results.
	fmt.Println("\n======================================")
	fmt.Println("           FINAL NODE STATUS")
	fmt.Println("======================================")

	for _, node := range nodes {
		fmt.Printf("\nComputer: %s (%s)\n", node.Name, node.City)

		if node.HeartbeatConfirmed {
			fmt.Println("Monitoring result: ALIVE (heartbeat confirmed)")
		} else {
			fmt.Println("Monitoring result: SUSPECTED UNAVAILABLE")
			fmt.Println("Reason: Heartbeat timeout exceeded")
		}
	}

	fmt.Println("\n======================================")
	fmt.Println("       MONITORING COMPLETED")
	fmt.Println("======================================")
}
