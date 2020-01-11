package latest

import (
	// openapi "github.com/go-openapi/spec"
	corev1 "k8s.io/api/core/v1"
)

// Manifest stores a kuda manifest.
type Manifest struct {
	ManivestVersion string `yaml:"kudaManifestVersion"`
	Version         string `yaml:"version,omitempty"`
	Name            string `yaml:"name"`
	Meta            Meta   `yaml:"meta,omitempty"`
	Deploy          Config `yaml:"deploy"`
	Dev             Config `yaml:"dev,omitempty"`
	// Paths           *openapi.Paths `yaml:"paths,omitempty"`
}

// Meta stores the metadata.
type Meta struct {
	Author     string `yaml:"author,omitempty"`
	Repository string `yaml:"repository,omitempty"`
	License    string `yaml:"license,omitempty"`
}

// Config stores a deployment config.
type Config struct {
	Dockerfile string          `yaml:"dockerfile,omitempty"`
	Entrypoint Entrypoint      `yaml:"entrypoint,omitempty"`
	Sync       []string        `yaml:"sync,omitempty"`
	Env        []corev1.EnvVar `yaml:"env,omitempty"`
}

// Entrypoint stores an API entrypoint.
type Entrypoint struct {
	Command string   `yaml:"command"`
	Args    []string `yaml:"args,omitempty"`
}
