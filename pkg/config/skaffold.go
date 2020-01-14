package config

import (
	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	latest "github.com/cyrildiagne/kuda/pkg/manifest/latest"
)

// GenerateSkaffoldConfig generate skaffold yaml specifics to the Kuda workflow.
func GenerateSkaffoldConfig(service ServiceSummary, manifest latest.Config, knativeFile string) (v1.SkaffoldConfig, error) {

	var sync *v1.Sync
	if manifest.Sync != nil {
		sync = &v1.Sync{
			Manual: []*v1.SyncRule{},
		}
		for _, s := range manifest.Sync {
			sync.Manual = append(sync.Manual, &v1.SyncRule{Src: s, Dest: "."})
		}
	}

	artifact := v1.Artifact{
		// The endpoint image name.
		ImageName: service.DockerArtifact,
		// Which Dockerfile to build.
		ArtifactType: v1.ArtifactType{
			DockerArtifact: &v1.DockerArtifact{
				DockerfilePath: manifest.Dockerfile,
			},
		},
		// Sync rules.
		Sync: sync,
	}

	tagPolicy := v1.TagPolicy{
		EnvTemplateTagger: &v1.EnvTemplateTagger{
			Template: "{{.IMAGE_NAME}}:{{.API_VERSION}}",
		},
	}

	build := v1.BuildConfig{
		Artifacts: []*v1.Artifact{&artifact},
		BuildType: *service.BuildType,
		TagPolicy: tagPolicy,
	}

	deploy := v1.DeployConfig{
		DeployType: v1.DeployType{
			// Location of the manifest file
			KubectlDeploy: &v1.KubectlDeploy{
				Manifests: []string{knativeFile},
			},
		},
	}

	config := v1.SkaffoldConfig{
		APIVersion: v1.Version,
		Kind:       "Config",
		Pipeline: v1.Pipeline{
			Build:  build,
			Deploy: deploy,
		},
	}

	return config, nil
}
