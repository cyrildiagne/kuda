package config

import (
	"testing"

	config "github.com/cyrildiagne/kuda/pkg/kuda/config"
	"gotest.tools/assert"
)

func TestGetDockerfileArtifactName(t *testing.T) {
	testUserCfg := config.UserConfig{
		Deployer: config.DeployerType{
			Skaffold: &config.SkaffoldDeployerConfig{
				DockerRegistry: "test-registry",
			},
		},
	}
	// Test without extension.
	name := GetDockerfileArtifactName(testUserCfg, "test")
	assert.Equal(t, "test-registry/test", name)

	// Test with "-dev" extension which should be removed in docker image name.
	name = GetDockerfileArtifactName(testUserCfg, "test-dev")
	assert.Equal(t, "test-registry/test", name)
}
