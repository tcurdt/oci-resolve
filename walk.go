package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

func processFiles(inputPath string, outputPath string) error {
	files, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		if file.IsDir() {
			processFiles(fmt.Sprintf("%s/%s", inputPath, file.Name()), outputPath)
			continue
		}

		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			filePath := fmt.Sprintf("%s/%s", inputPath, file.Name())
			err := processFile(filePath, outputPath)
			if err != nil {
				fmt.Printf("error processing file %s: %v\n", filePath, err)
			}
		}
	}

	return nil
}

func processFile(filePath string, outputPath string) error {

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(strings.NewReader(string(content)))

	for {
		var doc yaml.Node
		err := decoder.Decode(&doc)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Printf("error decoding YAML: %v\n", err)
			return err
		}
		processNode(&doc)
	}

	// var yamlDoc yaml.Node
	// err = yaml.Unmarshal(content, &yamlDoc)
	// if err != nil {
	// 	return err
	// }

	// updatedContent, err := yaml.Marshal(&yamlDoc)
	// if err != nil {
	// 	return err
	// }

	// err = ioutil.WriteFile(filePath, updatedContent, 0644)
	// if err != nil {
	// 	return err
	// }

	return nil
}
