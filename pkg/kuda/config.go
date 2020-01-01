package kuda

import (
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
)

// Config stores the general config.
type Config struct {
	Name      string
	Namespace string
	// Local config
	Dockerfile string
	// File Generation config
	ConfigFolder string
	KserviceFile string
	SkaffoldFile string
	// Artifacts
	DockerDestImage string
	// Dev config
	Sync    []string
	Command string
	Args    []string
	Env     []corev1.EnvVar
}

// Manifests contains the yamls of the kservice and skaffold manifests.
type Manifests struct {
	Kservice string
	Skaffold string
}

// AddDevConfigFlask initializes Config Dev props for Flask projects.
func (cfg *Config) AddDevConfigFlask() {
	cfg.Sync = []string{"**/*.py"}
	cfg.Command = "python3"
	cfg.Args = []string{"main.py"}
	cfg.Env = []corev1.EnvVar{
		{Name: "FLASK_ENV", Value: "development"},
	}
}

// SetFilesSuffix initializes Config Dev props for Flask projects.
func (cfg *Config) SetFilesSuffix(suffix string) {
	kExt := filepath.Ext(cfg.KserviceFile)
	cfg.KserviceFile = cfg.KserviceFile[:len(cfg.KserviceFile)-len(kExt)] + suffix + kExt
	cfg.SkaffoldFile = cfg.SkaffoldFile[:len(cfg.SkaffoldFile)-len(kExt)] + suffix + kExt
}

// NewConfig create new instance of Config.
func NewConfig(name string, namespace string) Config {
	cfg := Config{}
	cfg.Name = name
	cfg.Namespace = namespace
	cfg.ConfigFolder = "./.kuda"
	cfg.Dockerfile = "./Dockerfile"
	cfg.KserviceFile = cfg.ConfigFolder + "/" + "service.yml"
	cfg.SkaffoldFile = cfg.ConfigFolder + "/" + "skaffold.yml"
	return cfg
}

// IsValid checks wether or not a config is valid.
func (cfg *Config) IsValid() bool {
	if cfg.Name == "" || cfg.Namespace == "" {
		return false
	}
	if cfg.ConfigFolder == "" {
		return false
	}
	if cfg.Dockerfile == "" {
		return false
	}
	if cfg.KserviceFile == "" || cfg.SkaffoldFile == "" {
		return false
	}
	return true
}

// GenerateConfigFiles creates the configuration files from a fully qualitfied url.
func GenerateConfigFiles(cfg Config) (*Manifests, error) {

	// Write the service manifest.
	knativeManifest, err := GenerateKnativeConfigYAML(cfg)
	if err != nil {
		return nil, err
	}
	// Write the skaffold manifest.
	skaffoldManifest, err := GenerateSkaffoldConfigYAML(cfg)
	if err != nil {
		return nil, err
	}

	return &Manifests{
		Kservice: knativeManifest,
		Skaffold: skaffoldManifest,
	}, nil
}
