package config

import (
	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/auth"
)

// UserConfig stores a local user configuration.
type UserConfig struct {
	Namespace string       `yaml:"namespace"`
	Deployer  DeployerType `yaml:"deployer"`
}

// DeployerType stores the deployers configs.
type DeployerType struct {
	Skaffold *SkaffoldDeployerConfig `yaml:",omitempty"`
	Remote   *RemoteDeployerConfig   `yaml:",omitempty"`
}

// SkaffoldDeployerConfig stores a skaffold deployer config.
type SkaffoldDeployerConfig struct {
	// The destination Docker Registry. eg: gcr.io/project-name.
	DockerRegistry string `yaml:"dockerRegistry"`
	// Where the manifests should be written.
	ConfigFolder string `yaml:"configFolder,omitempty"`
}

// RemoteDeployerConfig stores a remote deployer config.
type RemoteDeployerConfig struct {
	AuthURL     string     `yaml:"auth"`
	DeployerURL string     `yaml:"deployer"`
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
