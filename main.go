package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"
	yaml "gopkg.in/yaml.v3"
)

type Image struct {
	Name string
	Tag  string
}

func (i Image) String() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

type Registry struct {
	Name string `yaml:"name"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type ResolveCmd struct {
	RegistryFiles []string `arg:"-r,--registry,separate" help:"registry file(s)"`
	Images        []string `arg:"-i,--image,separate" help:"image in name:tag format"`
	Output        string   `arg:"-o" help:"Output format (e.g., kustomize, json)" default:"kustomize"`
}

var args struct {
	Resolve *ResolveCmd `arg:"subcommand:resolve"`
	Verbose bool        `arg:"--verbose" help:"verbose output" default:"false"`
}

func parseRegistries(filenames []string) ([]Registry, error) {
	var registries []Registry

	for _, filename := range filenames {

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return []Registry{}, err
		}

		var registry Registry
		err = yaml.Unmarshal(data, &registry)
		if err != nil {
			return []Registry{}, err
		}

		registries = append(registries, registry)
	}
	return registries, nil
}

func parseImages(values []string) ([]Image, error) {
	var images []Image

	for _, value := range values {

		parts := strings.SplitN(value, ":", 2)

		if len(parts) != 2 {
			return []Image{}, fmt.Errorf("invalid image format: %s", value)
		}

		name := parts[0]
		tag := parts[1]

		images = append(images, Image{Name: name, Tag: tag})
	}
	return images, nil
}

func kustomize(resolved map[Image]string) {
	for image, sha := range resolved {
		fmt.Printf("kustomize edit set image %s=%s:%s\n", image.String(), image.Name, sha)
	}
}

func main() {

	arg.MustParse(&args)

	switch {
	case args.Resolve != nil:
		registries, err := parseRegistries(args.Resolve.RegistryFiles)
		if err != nil {
			LogError("error parsing registries: %v", err)
			os.Exit(1)
		}
		images, err := parseImages(args.Resolve.Images)
		if err != nil {
			LogError("error parsing images: %v", err)
			os.Exit(1)
		}
		LogInfo("registries: %v", registries)
		LogInfo("images: %v", images)
		LogInfo("output: %s", args.Resolve.Output)

		resolved, err := resolveImages(registries, images)
		if err != nil {
			LogError("error resolving images: %v", err)
			os.Exit(1)
		}
		if args.Resolve.Output == "kustomize" {
			kustomize(resolved)
			return
		}

		LogError("unknown output format", args.Resolve.Output)
		os.Exit(1)

	default:
		LogError("unknown command")
		os.Exit(1)
	}
}
