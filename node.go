package main

import "fmt"

type Node struct {
    Name string
    City string
    CPU  string
}

func main() {
    node := Node{
        Name: "computer-1",
        City: "Dhaka",
        CPU:  "AMD",
    }

    fmt.Println(node.Name)
    fmt.Println(node.City)
    fmt.Println(node.CPU)
}