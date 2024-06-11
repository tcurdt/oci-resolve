package main

import (
	yaml "gopkg.in/yaml.v3"
)

type Consumer interface {
	Consume(filePath string, doc *yaml.Node) error
}

func extractKind(doc *yaml.Node) string {
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		root := doc.Content[0]
		if root.Kind == yaml.MappingNode {
			for i := 0; i < len(root.Content); i += 2 {
				keyNode := root.Content[i]
				valueNode := root.Content[i+1]
				if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "kind" {
					return valueNode.Value
				}
			}
		}
	}
	return ""
}
