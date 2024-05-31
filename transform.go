package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func processNode(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode:
		processNode(node.Content[0])
	case yaml.MappingNode:
		processMappingNodeWithComments(node)
	case yaml.ScalarNode:
	case yaml.SequenceNode:
		for _, contentNode := range node.Content {
			processNode(contentNode)
		}
	default:
		fmt.Println("unexpected node kind:", node.Kind)
	}

}

func processMappingNodeWithComments(node *yaml.Node) {
	if node == nil || node.Kind != yaml.MappingNode {
		return
	}

	if len(node.Content) > 1 {
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "image" {

				fmt.Printf("found image: %s (%s)\n", valueNode.Value, valueNode.LineComment)

				if len(valueNode.LineComment) > 0 {

					// comment := string(valueNode.LineComment)

					// tag := extractTagFromComment(comment)
					// if tag != "" {
					// 	fmt.Println("found tag:", tag)
					// 	valueNode.Value = modifyImageValue(valueNode.Value, tag)
					// }

				}
			}

			processNode(valueNode)
		}
	}
}
