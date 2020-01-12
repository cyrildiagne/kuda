package config

import (
	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/auth"
)

// UserConfig stores a local user configuration.
type UserConfig struct {
	Namespace string         `yaml:"namespace"`
	Provider  ProviderConfig `yaml:"provider"`
}

// ProviderConfig stores a remote deployer config.
type ProviderConfig struct {
	AuthURL     string     `yaml:"auth"`
	ApiURL string     `yaml:"deployer"`
	User        *auth.User `yaml:"user"`
}

// Helpers

// ServiceSummary stores a summary of a knative service.
type ServiceSummary struct {
	Name           string
	Namespace      string
	DockerArtifact string
	BuildType      v1.BuildType
}
