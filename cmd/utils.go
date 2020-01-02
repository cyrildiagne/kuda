package cmd

import (
	"io/ioutil"
	"os"
	"reflect"

	"github.com/cyrildiagne/kuda/pkg/kuda/manifest/latest"
	yaml "gopkg.in/yaml.v2"
)

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
