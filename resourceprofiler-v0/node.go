package main

import "errors"

type Node struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	GPU      string `json:"gpu"`
	VRAM     int    `json:"vram_gb"`
}

func validateNode(node Node) error {
	if node.Name == "" {
		return errors.New("node name is required")
	}

	if node.GPU == "" {
		return errors.New("GPU model is required")
	}

	if node.VRAM <= 0 {
		return errors.New("VRAM must be greater than zero")
	}

	return nil
}
