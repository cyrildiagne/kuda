package kuda

import (
	"testing"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"gotest.tools/assert"
)

func TestGenerateSkaffoldConfig(t *testing.T) {

	config := GetTestConfig()

	result, err := GenerateSkaffoldConfig(config)
	if err != nil {
		t.Errorf("err")
	}

	CheckDeepEqual(t, result.APIVersion, v1.Version)
	CheckDeepEqual(t, result.Kind, "Config")

	artifacts := []*v1.Artifact{
		{
			ImageName: "test.io/test/test",
			ArtifactType: v1.ArtifactType{
				DockerArtifact: &v1.DockerArtifact{
					DockerfilePath: "./Dockerfile",
				},
			},
		},
	}
	CheckDeepEqual(t, result.Pipeline.Build.Artifacts, artifacts)

	deploy := v1.DeployConfig{
		DeployType: v1.DeployType{
			KubectlDeploy: &v1.KubectlDeploy{
				Manifests: []string{"./.kuda/service.yml"},
			},
		},
	}
	CheckDeepEqual(t, result.Deploy, deploy)
}

func TestGenerateSkaffoldDevConfig(t *testing.T) {

	emptyConfig := Config{}
	_, e := GenerateSkaffoldConfig(emptyConfig)
	assert.Error(t, e, "invalid config")

	config := GetTestDevConfig()
	result, err := GenerateSkaffoldConfig(config)
	if err != nil {
		t.Errorf("err")
	}

	artifacts := []*v1.Artifact{
		{
			ImageName: "test.io/test/test",
			ArtifactType: v1.ArtifactType{
				DockerArtifact: &v1.DockerArtifact{
					DockerfilePath: "./Dockerfile",
				},
			},
			Sync: &v1.Sync{
				Manual: []*v1.SyncRule{{Src: "**/*.py", Dest: "."}},
			},
		},
	}
	CheckDeepEqual(t, result.Pipeline.Build.Artifacts, artifacts)
}
