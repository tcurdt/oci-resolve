package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type HashTransformer struct{}

// kind, namespace, name

func NewHashTransformer() *HashTransformer {
	return &HashTransformer{}
}

func (c *HashTransformer) Consume(filePath string, doc *yaml.Node) error {
	return c.processNode(doc, []string{})
}

func (c *HashTransformer) processNode(node *yaml.Node, path []string) error {
	if node == nil {
		return nil
	}

	switch node.Kind {
	case yaml.DocumentNode:
		err := c.processNode(node.Content[0], []string{})
		if err != nil {
			return err
		}
	case yaml.MappingNode:
		err := c.processMappingNode(node, append(path, node.Value))
		if err != nil {
			return err
		}
	case yaml.ScalarNode:
	case yaml.SequenceNode:
		for i, contentNode := range node.Content {
			err := c.processNode(contentNode, append(path, fmt.Sprintf("[%d]", i)))
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unexpected node kind: %v", node.Kind)
	}

	return nil
}

func (c *HashTransformer) processMappingNode(node *yaml.Node, path []string) error {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	if len(node.Content) > 1 {
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			fullPath := append(path, keyNode.Value)

			_ = fullPath

			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "name" {
				// LogDebug("found %v: %s\n", fullPath, valueNode.Value)
			}

			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "app" {
				// LogDebug("found %v: %s\n", fullPath, valueNode.Value)
			}

			err := c.processNode(valueNode, append(path, keyNode.Value))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
