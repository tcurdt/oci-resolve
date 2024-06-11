package main

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type WriteTransformer struct {
	prefix        string
	filePathsSeen map[string]bool
}

func NewWriteTransformer() *WriteTransformer {
	return &WriteTransformer{
		prefix:        "",
		filePathsSeen: make(map[string]bool),
	}
}

func (c *WriteTransformer) Consume(filePath string, doc *yaml.Node) error {

	seen := c.filePathsSeen[filePath]

	var flags int
	if seen {
		flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	} else {
		flags = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}

	f, err := os.OpenFile(filepath.Join(c.prefix, filePath), flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// kind := extractKind(doc)
	// fmt.Printf("--- %s\n", kind)

	if seen {
		_, err := io.WriteString(f, "---\n")
		if err != nil {
			return err
		}
	}

	encoder := yaml.NewEncoder(f)
	defer encoder.Close()
	if err := encoder.Encode(doc); err != nil {
		return err
	}

	c.filePathsSeen[filePath] = true

	return nil
}
