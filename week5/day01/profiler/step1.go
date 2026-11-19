package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("architecture:", runtime.GOARCH)
}
