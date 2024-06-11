package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type PrintTransformer struct{}

func NewPrintTransformer() *PrintTransformer {
	return &PrintTransformer{}
}

func (c *PrintTransformer) Consume(filePath string, doc *yaml.Node) error {

	// kind := extractKind(doc)
	// fmt.Printf("--- %s\n", kind)

	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	if err := encoder.Encode(doc); err != nil {
		return err
	}

	fmt.Println("---")

	return nil
}
