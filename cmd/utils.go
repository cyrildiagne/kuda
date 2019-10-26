/*
Package cmd -

Copyright © 2019 Cyril Diagne <diagne.cyril@gmail.com>.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/viper"
)

// GetImage retrieves the provider docker image to use based on the user setting
// And the version compatible with the compiled CLI.
func GetImage() (string, error) {
	// Get provider from config.
	provider := viper.GetString("provider")
	var version string
	if provider == "gcp" {
		version = gcpProviderVersion
	} else if provider == "aws" {
		version = awsProviderVersion
	} else {
		return "", errors.New("provider unknown")
	}
	// Setup the provider's image.
	image := "gcr.io/kuda-project/provider-" + provider + ":" + version
	return image, nil
}

// RunProviderCommand runs the provider docker image and retrieves local
// environment variables that match a provider id.
func RunProviderCommand(opts docker.CommandOption) error {
	// Retrieve the provider image.
	image, err := GetImage()
	if err != nil {
		return errors.New("could not retrieve provider image")
	}
	opts.Image = image
	// Environment variables for the Docker image.
	// Convert all the viper configs to environment variable
	// in the format KUDA_* where * is the config uppercased.
	for k, e := range viper.AllSettings() {
		key := "KUDA_" + strings.ToUpper(k)
		value := fmt.Sprintf("%v", e)
		opts.AppendEnv = append(opts.AppendEnv, key+"="+value)
	}
	return docker.RunDockerCommand(opts)
}
