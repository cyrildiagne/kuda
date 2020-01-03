package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	skaffoldCfg "github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	yaml "gopkg.in/yaml.v2"
)

func generateSkaffoldConfigFiles(config latest.Config, name string, folder string) (string, error) {
	// Make sure output folder exists.
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0700)
	}

	// Generate the knative yaml file.
	knativeCfg, err := skaffoldCfg.GenerateKnativeConfig(name, config, cfg)
	if err != nil {
		return "", err
	}
	knativeYAML, err := skaffoldCfg.MarshalKnativeConfig(knativeCfg)
	if err != nil {
		return "", err
	}
	knativeFile := filepath.FromSlash(folder + "/knative-" + name + ".yaml")
	if err := writeYAML(knativeYAML, knativeFile); err != nil {
		return "", err
	}

	// Generate the skaffold yaml file.
	skaffoldCfg, err := skaffoldCfg.GenerateSkaffoldConfig(name, config, cfg, knativeFile)
	if err != nil {
		return "", err
	}
	skaffoldYAML, err := yaml.Marshal(skaffoldCfg)
	if err != nil {
		return "", err
	}
	skaffoldFile := filepath.FromSlash(folder + "/skaffold-" + name + ".yaml")
	if err := writeYAML(skaffoldYAML, skaffoldFile); err != nil {
		return "", err
	}

	return skaffoldFile, nil
}

func loadManifest(manifestFile string, manifest *latest.Manifest) error {
	// Check if file exists.
	if _, err := os.Stat(manifestFile); err != nil {
		return err
	}
	// Load file.
	data, err := ioutil.ReadFile(manifestFile)
	if err != nil {
		return err
	}
	// Parse.
	if err = yaml.Unmarshal(data, manifest); err != nil {
		return err
	}
	// Dev extends values from Deploy.
	valuesDev := reflect.ValueOf(&manifest.Dev).Elem()
	valuesDeploy := reflect.ValueOf(&manifest.Deploy).Elem()
	for i := 0; i < valuesDev.NumField(); i++ {
		fDev := valuesDev.Field(i)
		// Check if the field is zero-valued.
		if reflect.DeepEqual(fDev.Interface(), reflect.Zero(fDev.Type()).Interface()) {
			fDeploy := valuesDeploy.Field(i)
			fDev.Set(fDeploy)
		}
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
