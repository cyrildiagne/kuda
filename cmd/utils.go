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
	"fmt"
	"strings"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/viper"
)

// RunDockerWithProviderEnvs retrieves local environment variables
// that match a provider id and runs a docker image.
func RunDockerWithProviderEnvs(opts docker.CommandOption) error {
	// Environment variables for the Docker image.
	// We look for all the configs that start with the provider name
	// "gcp" and convert them in the environment variable
	// format KUDA_GCP_*
	provider := viper.GetString("provider")
	for k, e := range viper.AllSettings() {
		if strings.HasPrefix(k, provider) {
			key := "KUDA_" + strings.ToUpper(k)
			value := fmt.Sprintf("%v", e)
			opts.AppendEnv = append(opts.AppendEnv, key+"="+value)
		}
	}
	return docker.RunDockerCommand(opts)
}
