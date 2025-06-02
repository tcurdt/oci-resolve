package main

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type ResolveTransformer struct {
	registries []Registry
	images     []Image
}

func NewResolveTransformer(registries []Registry, images []Image) *ResolveTransformer {
	return &ResolveTransformer{registries, images}
}

func (c *ResolveTransformer) Consume(filePath string, doc *yaml.Node) error {
	return c.processNode(doc, []string{})
}

func (c *ResolveTransformer) processNode(node *yaml.Node, path []string) error {
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

func (c *ResolveTransformer) processMappingNode(node *yaml.Node, path []string) error {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	if len(node.Content) > 1 {
		for i := 0; i < len(node.Content)-1; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			fullPath := append(path, keyNode.Value)

			if keyNode.Kind == yaml.ScalarNode && keyNode.Value == "image" {

				// fmt.Printf("found image %v %s (%s)", fullPath, valueNode.Value, valueNode.LineComment)

				// check if matches images
				// check if annotation

				if c.hasCommentTag(valueNode) {
					LogDebug("found image %v %s (%s)", fullPath, valueNode.Value, valueNode.LineComment)
					image, err := parseImage(valueNode.Value)
					if err != nil {
						return err
					}
					resolved, err := resolveImage(c.registries, image)
					if err != nil {
						LogError("failed to resolve image [%v] err=%v", image.String(), err)
						return nil
						// return err
					}
					LogInfo("resolved image [%v] => %s", image.String(), resolved)
					valueNode.Value = Image{image.Name, resolved}.String()
					return nil
				}
			}

			err := c.processNode(valueNode, append(path, keyNode.Value))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *ResolveTransformer) hasCommentTag(valueNode *yaml.Node) bool {

	if len(valueNode.LineComment) > 0 {

		comment := string(valueNode.LineComment)

		if strings.Contains(comment, "resolve") {
			return true
		}
	}

	return false
}
