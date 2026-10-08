package main

import "fmt"

type Node struct {
	Name string
	CPU  string
}

func main() {
	node1 := Node{
		Name: "Node 1",
		CPU:  "AMD",
	}

	node2 := Node{
		Name: "Node 2",
		CPU:  "Intel",
	}

	node3 := Node{
		Name: "Node 3",
		CPU:  "AMD",
	}

	fmt.Println(node1)
	fmt.Println(node2)
	fmt.Println(node3)
}