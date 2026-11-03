package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Sender   string
	Receiver string
	Content  string
}

type Capability struct {
	Name     string `json:"name"`
	GPU      string `json:"gpu"`
	VRAMGB   int    `json:"vram_gb"`
	Location string `json:"location"`
}

func main() {
	capabilityA := Capability{
		Name:     "Computer A",
		GPU:      "RTX 5060",
		VRAMGB:   8,
		Location: "Dhaka",
	}

	fmt.Println("Computer A Hardware:")
	capabilityJSON, err := json.Marshal(capabilityA)

	if err != nil {
		fmt.Println("JSON encoding failed:", err)
		return
	}

	fmt.Println("Capability JSON:")
	fmt.Println(string(capabilityJSON))

	message := Message{
		Sender:   "Computer A",
		Receiver: "Computer B",
		Content:  string(capabilityJSON),
	}

	fmt.Println("MESSAGE TRACE:")
	fmt.Println(message.Sender, "---->", message.Receiver)

	fmt.Println("Sending message:")
	fmt.Println("Sender:", message.Sender)
	fmt.Println("Receiver:", message.Receiver)
	fmt.Println("Content:", message.Content)
	var receivedCapability Capability

	err = json.Unmarshal([]byte(message.Content), &receivedCapability)

	if err != nil {
		fmt.Println("Receiving failed:", err)
		return
	}

	fmt.Println()
	fmt.Println("Computer B received hardware information:")
	fmt.Println("Name:", receivedCapability.Name)
	fmt.Println("GPU:", receivedCapability.GPU)
	fmt.Println("VRAM:", receivedCapability.VRAMGB, "GB")
	fmt.Println("Location:", receivedCapability.Location)
}
