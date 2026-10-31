package main

import "testing"

func TestValidateNode(t *testing.T) {
    node := Node{
        Name: "computer-1",
        GPU:  "RTX 5060",
        VRAM: 8,
    }

    err := validateNode(node)

    if err != nil {
        t.Fatalf("Expected valid node, got error: %v", err)
    }
}
func TestValidateNodeRejectsZeroVRAM(t *testing.T) {
    node := Node{
        Name: "computer-2",
        GPU:  "RTX 5070",
        VRAM: 0,
    }

    err := validateNode(node)

    if err == nil {
        t.Fatal("Expected an error for zero VRAM, got nil")
    }
}

func TestValidateNodeRejectsMissingGPU(t *testing.T) {
    node := Node{
        Name: "computer-3",
        GPU:  "",
        VRAM: 8,
    }

    err := validateNode(node)

    if err == nil {
        t.Fatal("Expected an error for missing GPU, got nil")
    }
}

func TestValidateNodeRejectsMissingName(t *testing.T) {
    node := Node{
        Name: "",
        GPU:  "RTX 5060",
        VRAM: 8,
    }

    err := validateNode(node)

    if err == nil {
        t.Fatal("Expected an error for missing name, got nil")
    }
}
