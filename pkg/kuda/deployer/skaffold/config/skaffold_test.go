package config

import (
	"testing"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	config "github.com/cyrildiagne/kuda/pkg/kuda/config"
	latest "github.com/cyrildiagne/kuda/pkg/kuda/manifest/latest"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
)

func TestGenerateSkaffoldConfig(t *testing.T) {

	name := "test-name"
	cfg := latest.Config{
		Dockerfile: "test-file",
		Sync:       []string{"test-sync"},
	}
	userCfg := config.UserConfig{
		Namespace: "test-namespace",
		Deployer: config.DeployerType{
			Skaffold: &config.SkaffoldDeployerConfig{
				DockerRegistry: "test-registry",
			},
		},
	}
	knativeConfig := "test-knative-file"

	result, err := GenerateSkaffoldConfig(name, cfg, userCfg, knativeConfig)
	if err != nil {
		t.Errorf("err")
	}

	assert.Equal(t, result.APIVersion, v1.Version)
	assert.Equal(t, result.Kind, "Config")

	artifacts := []*v1.Artifact{
		{
			ImageName: "test-registry/test-name",
			ArtifactType: v1.ArtifactType{
				DockerArtifact: &v1.DockerArtifact{
					DockerfilePath: "test-file",
				},
			},
			Sync: &v1.Sync{
				Manual: []*v1.SyncRule{{Src: "test-sync", Dest: "."}},
			},
		},
	}
	if diff := cmp.Diff(result.Pipeline.Build.Artifacts, artifacts); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}

	deploy := v1.DeployConfig{
		DeployType: v1.DeployType{
			KubectlDeploy: &v1.KubectlDeploy{
				Manifests: []string{"test-knative-file"},
			},
		},
	}
	if diff := cmp.Diff(result.Deploy, deploy); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
