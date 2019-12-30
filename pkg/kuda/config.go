package kuda

import (
	"net/url"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// Config stores the general config.
type Config struct {
	URLConfig URLConfig
	// Artifacts
	DockerDestImage string
	ManifestFile    string
	// Local config
	Dockerfile   string
	ConfigFolder string
	// Dev config
	DevConfig *DevConfig
}

// DevConfig contains dev overrides.
type DevConfig struct {
	Sync    []string
	Command string
	Args    []string
	Env     []corev1.EnvVar
}

// URLConfig stores the API config extracted from the url.
type URLConfig struct {
	Scheme    string
	Name      string
	Namespace string
	Domain    string
}

// NewDevConfigFlask create new instance of DevConfig for Flask projects.
func NewDevConfigFlask() DevConfig {
	cfg := DevConfig{}
	cfg.Sync = []string{"**/*.py"}
	cfg.Command = "python3"
	cfg.Args = []string{"main.py"}
	cfg.Env = []corev1.EnvVar{
		{Name: "FLASK_ENV", Value: "development"},
	}
	return cfg
}

// NewConfig create new instance of Config.
func NewConfig() Config {
	cfg := Config{}
	cfg.ConfigFolder = "./.kuda"
	cfg.Dockerfile = "./Dockerfile"
	return cfg
}

// ParseURL extract the config from an URL.
func ParseURL(dURL string) (*URLConfig, error) {
	u, err := url.Parse(dURL)
	if err != nil {
		return nil, err
	}

	split := strings.Split(u.Host, ".")

	config := URLConfig{
		Scheme:    u.Scheme,
		Name:      split[0],
		Namespace: split[1],
		Domain:    strings.Join(split[2:], "."),
	}

	return &config, nil
}

// WriteManifest writes the manifest to disk
func WriteManifest(content string, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// GenerateConfigFiles creates the configuration files from a fully qualitfied url.
func GenerateConfigFiles(dURL string, dockerRegistry string) error {

	// Infer config from URL.
	urlCfg, err := ParseURL(dURL)
	if err != nil {
		return err
	}

	// Docker config.
	dockerDestImage := dockerRegistry + "/" + urlCfg.Name

	cfg := NewConfig()

	// Make sure config folder exists.
	if _, err := os.Stat(cfg.ConfigFolder); os.IsNotExist(err) {
		os.Mkdir(cfg.ConfigFolder, 0700)
	}

	// Create config.
	cfg.URLConfig = *urlCfg
	cfg.DockerDestImage = dockerDestImage

	// -- PROD --

	cfg.ManifestFile = cfg.ConfigFolder + "/" + "service.yml"
	// Write the service manifest.
	knativeManifest, err := GenerateKnativeConfigYAML(cfg)
	if err != nil {
		return err
	}
	if err := WriteManifest(knativeManifest, cfg.ManifestFile); err != nil {
		return err
	}
	// Write the skaffold manifest.
	skaffoldManifest, err := GenerateSkaffoldConfigYAML(cfg)
	if err != nil {
		return err
	}
	if err := WriteManifest(skaffoldManifest, cfg.ConfigFolder+"/"+"knative.yml"); err != nil {
		return err
	}

	// -- DEV --

	cfg.ManifestFile = cfg.ConfigFolder + "/" + "service-dev.yml"
	devCfg := NewDevConfigFlask()
	cfg.DevConfig = &devCfg
	// Write the service manifest.
	knativeManifestDev, err := GenerateKnativeConfigYAML(cfg)
	if err != nil {
		return err
	}
	if err := WriteManifest(knativeManifestDev, cfg.ManifestFile); err != nil {
		return err
	}
	// Write the skaffold manifest.
	skaffoldManifestDev, err := GenerateSkaffoldConfigYAML(cfg)
	if err != nil {
		return err
	}
	if err := WriteManifest(skaffoldManifestDev, cfg.ConfigFolder+"/"+"knative-dev.yml"); err != nil {
		return err
	}

	return nil
}
