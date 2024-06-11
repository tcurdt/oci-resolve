package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func processFile(filename string, consumers []Consumer) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(content)))

	for {
		var doc yaml.Node
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		for _, consumer := range consumers {
			if err := consumer.Consume(filename, &doc); err != nil {
				return err
			}
		}
	}

	return nil
}

func traversePath(dir string, consumers []Consumer) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
			err := processFile(path, consumers)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
