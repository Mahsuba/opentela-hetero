package main

import "fmt"

type Node struct {
    Name     string
    GPUCount int
}

func validateNode(node Node) error {
    if node.Name == "" {
        return fmt.Errorf("node name cannot be empty")
    }

    if node.GPUCount < 0 {
        return fmt.Errorf("GPU count cannot be negative")
    }

    return nil
}

func main() {
    node := Node{
        Name: "computer-1",
        GPUCount: 4,
    }

    err := validateNode(node)

    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Node is valid!")
    }
}
