package latest

import (
	"errors"
	"io"
	"io/ioutil"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

// Load the content of a file into a manifest.
func (manifest *Manifest) Load(reader io.Reader) error {
	// Load manifest data.
	data, err := ioutil.ReadAll(reader)
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
	return manifest.CheckValid()
}

// CheckValid ensures a manifest fields are properly set.
func (manifest *Manifest) CheckValid() error {
	if manifest.Name == "" {
		return errors.New("name can't be empy")
	}
	return nil
}
