package config

import (
	"testing"

	"gotest.tools/assert"
)

func TestGetDockerfileArtifactName(t *testing.T) {
	testUserCfg := UserConfig{
		Deployer: DeployerType{
			Skaffold: &SkaffoldDeployerConfig{
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
