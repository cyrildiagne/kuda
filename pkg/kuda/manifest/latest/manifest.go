package latest

import (
	corev1 "k8s.io/api/core/v1"
)

// Manifest stores a kuda manifest.
type Manifest struct {
	Version string `yaml:"kudaManifestVersion"`
	Name    string `yaml:"name"`
	Deploy  Config `yaml:"deploy"`
	Dev     Config `yaml:"dev,omitempty"`
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
