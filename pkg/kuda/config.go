package kuda

import (
	"fmt"
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
	// Local config
	Dockerfile string
	// File Generation config
	ConfigFolder string
	KserviceFile string
	SkaffoldFile string
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
	Protocol  string
	Namespace string
	Name      string
	Domain    string
}

// Manifests contains the yamls of the kservice and skaffold manifests.
type Manifests struct {
	Config   Config
	Kservice string
	Skaffold string
}

// ConfigManifests contains the generated manifests.
type ConfigManifests struct {
	Prod Manifests
	Dev  Manifests
}

// NewURLConfig creates a new instance of URLConfig with default values.
func NewURLConfig() URLConfig {
	cfg := URLConfig{}
	cfg.Protocol = "http"
	cfg.Namespace = "default"
	return cfg
}

// NewDevConfigFlask creates a new instance of DevConfig for Flask projects.
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
	cfg.KserviceFile = cfg.ConfigFolder + "/" + "service.yml"
	cfg.SkaffoldFile = cfg.ConfigFolder + "/" + "skaffold.yml"
	return cfg
}

// ParseURL extract the config from an URL.
func ParseURL(dURL string) (*URLConfig, error) {
	u, err := url.Parse(dURL)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("Invalid URL %s \n"+
			"Should be a valid URL: {scheme}://{api-name}.{namespace}.{domain}", dURL)
	}

	split := strings.Split(u.Host, ".")

	config := NewURLConfig()
	if u.Scheme != "" {
		config.Protocol = u.Scheme
	}

	config.Name = split[0]
	config.Namespace = split[1]
	config.Domain = strings.Join(split[2:], ".")

	return &config, nil
}

// GenerateConfigFiles creates the configuration files from a fully qualitfied url.
func GenerateConfigFiles(dURL string, dockerRegistry string) (*ConfigManifests, error) {

	// Infer config from URL.
	urlCfg, err := ParseURL(dURL)
	if err != nil {
		return nil, err
	}

	// -- PROD --

	// Create config.
	cfg := NewConfig()
	cfg.URLConfig = *urlCfg
	cfg.DockerDestImage = dockerRegistry + "/" + urlCfg.Name

	// Make sure config folder exists.
	if _, err := os.Stat(cfg.ConfigFolder); os.IsNotExist(err) {
		os.Mkdir(cfg.ConfigFolder, 0700)
	}

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

	// -- DEV --

	// Create config.
	cfgDev := NewConfig()
	cfgDev.URLConfig = *urlCfg
	cfgDev.DockerDestImage = dockerRegistry + "/" + urlCfg.Name

	// Make sure config folder exists.
	if _, err := os.Stat(cfgDev.ConfigFolder); os.IsNotExist(err) {
		os.Mkdir(cfgDev.ConfigFolder, 0700)
	}

	devCfg := NewDevConfigFlask()
	cfgDev.DevConfig = &devCfg
	cfgDev.KserviceFile = cfgDev.ConfigFolder + "/" + "service-dev.yml"
	cfgDev.SkaffoldFile = cfgDev.ConfigFolder + "/" + "skaffold-dev.yml"
	// Write the service manifest.
	knativeManifestDev, err := GenerateKnativeConfigYAML(cfgDev)
	if err != nil {
		return nil, err
	}
	// Write the skaffold manifest.
	skaffoldManifestDev, err := GenerateSkaffoldConfigYAML(cfgDev)
	if err != nil {
		return nil, err
	}

	return &ConfigManifests{
		Prod: Manifests{
			Config:   cfg,
			Kservice: knativeManifest,
			Skaffold: skaffoldManifest,
		},
		Dev: Manifests{
			Config:   cfgDev,
			Kservice: knativeManifestDev,
			Skaffold: skaffoldManifestDev,
		},
	}, nil
}
