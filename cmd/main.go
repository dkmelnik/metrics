package main

import "fmt"

func main() {
	o := getOperations("")
	fmt.Println(len(o), o == nil)
}

func getOperations(id string) []float32 {
	operations := make([]float32, 0)
	if id == "" {
		return nil
	}
	// Add elements to operations
	return operations
}
