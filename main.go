package main

import (
	"fmt"
)

func main() {
	inputPath := "./in"
	outputPath := "./out"

	err := processFiles(inputPath, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
