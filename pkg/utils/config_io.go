package utils

import (
	"os"
	"path/filepath"

	config "github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	yaml "gopkg.in/yaml.v2"
)

// LoadManifest loads a manifest from disk.
func LoadManifest(manifestFile string) (*latest.Manifest, error) {
	// Ensure manifest exists
	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		return nil, err
	}
	// Open the file.
	manifestReader, err := os.Open(manifestFile)
	if err != nil {
		return nil, err
	}
	// Load the manifest.
	manifest := latest.Manifest{}
	if err := manifest.Load(manifestReader); err != nil {
		return nil, err
	}
	return &manifest, nil
}

// GenerateSkaffoldConfigFiles generates the skaffold config files to disk.
func GenerateSkaffoldConfigFiles(service config.ServiceSummary, appCfg latest.Config, folder string) error {
	// Make sure output folder exists.
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0700)
	}

	// Generate the knative yaml file.
	knativeCfg, err := config.GenerateKnativeConfig(service, appCfg)
	if err != nil {
		return err
	}
	knativeYAML, err := config.MarshalKnativeConfig(knativeCfg)
	if err != nil {
		return err
	}
	knativeFile := filepath.FromSlash(folder + "/knative.yaml")
	if err := writeYAML(knativeYAML, knativeFile); err != nil {
		return err
	}

	// Generate the skaffold yaml file.
	skaffoldCfg, err := config.GenerateSkaffoldConfig(service, appCfg, knativeFile)
	if err != nil {
		return err
	}
	skaffoldYAML, err := yaml.Marshal(skaffoldCfg)
	if err != nil {
		return err
	}
	skaffoldFile := filepath.FromSlash(folder + "/skaffold.yaml")
	if err := writeYAML(skaffoldYAML, skaffoldFile); err != nil {
		return err
	}

	return nil
}

func writeYAML(content []byte, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		return err
	}
	return nil
}
