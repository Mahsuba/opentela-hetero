
package main

import "fmt"

func main() {
    locations := map[string]string{
        "computer-1": "Dhaka",
        "computer-2": "London",
        "computer-3": "Zurich",
    }

    cpus := map[string]string{
        "computer-1": "AMD",
        "computer-2": "Intel",
        "computer-3": "AMD",
    }

    fmt.Println("Location:", locations["computer-2"])
    fmt.Println("CPU:", cpus["computer-2"])
    location, exists := locations["computer-99"]
    fmt.Println("Location:", location)
    fmt.Println("Does node exist?", exists)
}

