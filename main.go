package main

import (
	"fmt"
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

type TransformCmd struct {
	Paths         []string `arg:"positional" help:"path(s)"`
	RegistryFiles []string `arg:"-r,--registry,separate" help:"registry file(s)"`
	Images        []string `arg:"-i,--image,separate" help:"image in name:tag format"`
}

var args struct {
	Transform *TransformCmd `arg:"subcommand:transform"`
	LogLevel  string        `arg:"--log" help:"log level" default:"info"`
}

func parseRegistries(filenames []string) ([]Registry, error) {
	var registries []Registry

	for _, filename := range filenames {

		data, err := os.ReadFile(filename)
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

	name := os.Getenv("REGISTRY_NAME")
	if name != "" {
		user := os.Getenv("REGISTRY_USER")
		pass := os.Getenv("REGISTRY_PASS")
		registries = append(registries, Registry{
			Name: name,
			User: user,
			Pass: pass,
		})
	}

	return registries, nil
}

func parseImage(value string) (Image, error) {
	parts := strings.SplitN(value, ":", 2)

	if len(parts) != 2 {
		return Image{}, fmt.Errorf("invalid image format: %s", value)
	}

	name := parts[0]
	tag := parts[1]

	return Image{Name: name, Tag: tag}, nil
}

func parseImages(values []string) ([]Image, error) {
	var images []Image

	for _, value := range values {

		image, err := parseImage(value)
		if err != nil {
			return []Image{}, err
		}

		images = append(images, image)
	}

	return images, nil
}

func parseLogLevel(level string) (Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return Debug, nil
	case "info":
		return Info, nil
	case "warn":
		return Warn, nil
	case "error":
		return Error, nil
	default:
		return Info, fmt.Errorf("invalid log level: %s", level)
	}
}

func main() {

	arg.MustParse(&args)

	level, err := parseLogLevel(args.LogLevel)
	if err != nil {
		LogError("error parsing log level: %v", err)
		os.Exit(1)
	}
	SetLogLevel(level)

	switch {
	case args.Transform != nil:

		registries, err := parseRegistries(args.Transform.RegistryFiles)
		if err != nil {
			LogError("error parsing registries: %v", err)
			os.Exit(1)
		}

		images, err := parseImages(args.Transform.Images)
		if err != nil {
			LogError("error parsing images: %v", err)
			os.Exit(1)
		}

		LogDebug("registries: %v", registries)
		LogDebug("images: %v", images)

		consumers := []Consumer{
			NewResolveTransformer(registries, images),
			NewHashTransformer(),
			NewPrintTransformer(),
		}

		for _, path := range args.Transform.Paths {
			if err := traversePath(path, consumers); err != nil {
				LogError("error processing directory %q: %v\n", path, err)
				os.Exit(1)
			}
		}

	default:
		LogError("unknown command")
		os.Exit(1)
	}
}
