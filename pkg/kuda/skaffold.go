package kuda

import (
	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	yaml2 "gopkg.in/yaml.v2"
)

// GenerateSkaffoldConfigYAML generate yaml string.
func GenerateSkaffoldConfigYAML(cfg Config) (string, error) {
	config, err := GenerateSkaffoldConfig(cfg)
	if err != nil {
		return "", err
	}
	content, err := yaml2.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GenerateSkaffoldConfig generate skaffold yaml specifics to the Kuda workflow
// and based on the kuda.Config given as parameter.
func GenerateSkaffoldConfig(cfg Config) (v1.SkaffoldConfig, error) {

	var sync *v1.Sync
	if cfg.DevConfig != nil {
		sync = &v1.Sync{
			Manual: []*v1.SyncRule{},
		}
		for _, s := range cfg.DevConfig.Sync {
			sync.Manual = append(sync.Manual, &v1.SyncRule{Src: s, Dest: "."})
		}
	}

	artifact := v1.Artifact{
		// The endpoint image name.
		ImageName: cfg.DockerDestImage,
		// Which Dockerfile to build.
		ArtifactType: v1.ArtifactType{
			DockerArtifact: &v1.DockerArtifact{
				DockerfilePath: cfg.Dockerfile,
			},
		},
		// Sync rules.
		Sync: sync,
	}

	build := v1.BuildConfig{
		Artifacts: []*v1.Artifact{&artifact},
	}

	deploy := v1.DeployConfig{
		DeployType: v1.DeployType{
			// Location of the manifest file
			KubectlDeploy: &v1.KubectlDeploy{
				Manifests: []string{cfg.KserviceFile},
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
