package main

import (
	"context"
	"fmt"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/config"
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

// func extractTagFromComment(comment string) string {
// 	re := regexp.MustCompile(`# tag:(\w+)`)
// 	match := re.FindStringSubmatch(comment)
// 	if len(match) == 2 {
// 		return match[1]
// 	}
// 	return ""
// }

func or(a string, b string) string {
	if a != "" {
		return a
	}
	return b
}

// func resolveImages(registries []Registry, images []Image) (map[Image]string, error) {
// 	results := map[Image]string{}
// 	for _, image := range images {
// 		sha, err := resolveImage(registries, image)
// 		if err != nil {
// 			return results, err
// 		}
// 		results[image] = sha
// 	}
// 	return results, nil
// }

func resolveImage(registries []Registry, image Image) (string, error) {

	var hosts []config.Host
	for _, registry := range registries {
		host := config.Host{
			Name: registry.Name,
			User: registry.User,
			Pass: registry.Pass,
			// ReqPerSec: 100,
			// ReqConcurrent: 10,
		}

		// fmt.Printf("Name: %s\n", host.Name)
		// fmt.Printf("User: %s\n", host.User)
		// fmt.Printf("Pass: %s\n", host.Pass)

		hosts = append(hosts, host)
	}

	ctx := context.Background()

	// delayInit, _ := time.ParseDuration("0.05s")
	// delayMax, _ := time.ParseDuration("0.10s")

	rc := regclient.New(
		regclient.WithConfigHost(hosts...),
		// WithRetryDelay(delayInit, delayMax),
	)

	// regctl image inspect ghcr.io/aquasecurity/trivy:latest

	needle := image.String()

	LogDebug("resolving image [%s] (name=%s tag=%s)", needle, image.Name, image.Tag)

	r, err := ref.New(needle)
	if err != nil {
		return "", fmt.Errorf("failed creating getRef: %v", err)
	}

	m, err := rc.ManifestGet(ctx, r)
	if err != nil {
		return "", fmt.Errorf("failed running ManifestGet: %v", err)
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
	if err != nil {
		return "", fmt.Errorf("failed running ManifestGet: %v", err)
	}

	mi, ok := m.(manifest.Imager)
	if !ok {
		return "", fmt.Errorf("manifest does not support image methods")
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

	// LogDebug("SHA", imageConfig.Config.Labels["SHA"])
	// LogDebug("org.opencontainers.image.revision", imageConfig.Config.Labels["org.opencontainers.image.revision"])

	sha := or(
		imageConfig.Config.Labels["SHA"],
		imageConfig.Config.Labels["org.opencontainers.image.revision"],
	)
	if sha != "" {
		return "commit-" + sha, nil
	}

	return "", fmt.Errorf("could not resolve image %s", needle)
}
