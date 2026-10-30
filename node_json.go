package main

import (
	"encoding/json"
	"fmt"
)

type Node struct {
	Name string `json:"name"`
	City string `json:"city"`
	CPU  string `json:"cpu"`
}

func main() {
	node := Node{
		Name: "computer-1",
		City: "Dhaka",
		CPU:  "Intel",
	}

	data, err := json.MarshalIndent(node, "", "    ")

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Valid JSON:", json.Valid(data))
	fmt.Println(string(data))
}
