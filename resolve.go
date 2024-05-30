package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
	"github.com/regclient/regclient/types"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/platform"
	"github.com/regclient/regclient/types/ref"
)

//	type OCIConfig struct {
//		Architecture string `yaml:"architecture"`
//		Config       struct {
//			Env        []string `yaml:"Env"`
//			Entrypoint []string `yaml:"Entrypoint"`
//			Labels map[string]string `yaml:"Labels"`
//			OnBuild interface{} `yaml:"OnBuild"`
//		} `yaml:"config"`
//		Created time.Time `yaml:"created"`
//		History []struct {
//			Created    time.Time `yaml:"created"`
//			CreatedBy  string    `yaml:"created_by"`
//			EmptyLayer bool      `yaml:"empty_layer,omitempty"`
//			Comment    string    `yaml:"comment,omitempty"`
//		} `yaml:"history"`
//		Os     string `yaml:"os"`
//		Rootfs struct {
//			Type    string   `yaml:"type"`
//			DiffIds []string `yaml:"diff_ids"`
//		} `yaml:"rootfs"`
//	}

func extractTagFromComment(comment string) string {
	re := regexp.MustCompile(`# tag:(\w+)`)
	match := re.FindStringSubmatch(comment)
	if len(match) == 2 {
		return match[1]
	}
	return ""
}

func or(a string, b string) string {
	if a != "" {
		return a
	}
	return b
}

func resolve(image string, tag string) (string, error) {

	fmt.Printf("resolving image%s tag=%s\n", image, tag)

	hosts := []config.Host{
		{
			Name: "ghcr.io",
			User: "tcurdt",
			Pass: "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",
			// ReqPerSec: 100,
			// ReqConcurrent: 10,
		},
	}

	ctx := context.Background()

	// delayInit, _ := time.ParseDuration("0.05s")
	// delayMax, _ := time.ParseDuration("0.10s")

	rc := regclient.New(
		regclient.WithConfigHost(hosts...),
		// WithRetryDelay(delayInit, delayMax),
	)

	// regctl image inspect ghcr.io/aquasecurity/trivy:latest
	r, err := ref.New("ghcr.io/tcurdt/test-project:live")
	// r, err := ref.New("ghcr.io/aquasecurity/trivy:latest")
	if err != nil {
		return "", fmt.Errorf("Failed creating getRef: %v", err)
	}

	m, err := rc.ManifestGet(ctx, r)
	if err != nil {
		return "", fmt.Errorf("Failed running ManifestGet: %v", err)
	}

	// if manifest.GetMediaType(m) != types.MediaTypeDocker2Manifest {
	// 	// Unexpected media type: application/vnd.docker.distribution.manifest.list.v2+json
	// 	// return fmt.Errorf("Unexpected media type: %s", manifest.GetMediaType(m)).Error()
	// }

	// fmt.Printf("isList=%@\n", m.IsList())

	plat := platform.Local()
	desc, err := manifest.GetPlatformDesc(m, &plat)
	if err != nil {
		return "", err
	}

	m, err = rc.ManifestGet(ctx, r, regclient.WithManifestDesc(*desc))

	mi, ok := m.(manifest.Imager)
	if !ok {
		return "", fmt.Errorf("manifest does not support image methods%.0w", types.ErrUnsupportedMediaType)
	}

	cd, err := mi.GetConfig()
	if err != nil {
		return "", err
	}

	blobConfig, err := rc.BlobGetOCIConfig(ctx, r, cd)
	if err != nil {
		return "", err
	}

	imageConfig := blobConfig.GetConfig()

	// body, err := blobConfig.RawBody()
	// if err != nil {
	// 	return "", err
	// }

	// var imageConfig OCIConfig
	// err = yaml.Unmarshal(body, &imageConfig)
	// if err != nil {
	// 	return "", err
	// }

	sha := or(
		imageConfig.Config.Labels["org.opencontainers.image.revision"],
		imageConfig.Config.Labels["SHA"],
	)
	if sha != "" {
		fmt.Printf("sha = %s\n", sha)
		return "commit-" + sha, nil
	}

	// don't change anything
	fmt.Println("could not resolve. keeping as is")

	return tag, nil

	// var yamlDoc yaml.Node
	// err = yaml.Unmarshal(body, &yamlDoc)
	// if err != nil {
	// 	return "", err
	// }

	// var out strings.Builder
	// encoder := yaml.NewEncoder(&out)
	// encoder.SetIndent(2)
	// err = encoder.Encode(&yamlDoc)
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Println(out.String())
	// fmt.Printf("Out:\n%s\n", out.String())

	// var yamlBytes bytes.Buffer
	// yamlEncoder := yaml.NewEncoder(&yamlBytes)
	// yamlEncoder.SetIndent(2)
	// yamlEncoder.Encode(&yamlDoc)

	// yamlBytes, err := yaml.Marshal(&yamlDoc)
	// if err != nil {
	// 	return "", err
	// }
	// fmt.Printf("Doc:\n%s\n", yamlBytes)
}

func modifyImageValue(value string, tag string) string {

	parts := strings.SplitN(value, ":", 2)

	var image string
	if len(parts) > 0 {
		image = parts[0]
	} else {
		image = value
	}

	resolved, err := resolve(image, tag)
	if err != nil {
		fmt.Println(err)
		return value
	}

	return image + ":" + resolved
}
