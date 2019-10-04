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
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var provider = "gcp"
var project string
var credentials string

// gcpCmd represents the `setup gcp` command
var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Setup Kuda on GCP.",
	Long:  "Setup Kuda on GCP.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := setup(); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	setupCmd.AddCommand(gcpCmd)

	gcpCmd.PersistentFlags().StringVarP(&project, "project", "p", "",
		"GCP Project ID")
	gcpCmd.MarkPersistentFlagRequired("project")
	viper.BindPFlag("gcp_project_id", gcpCmd.PersistentFlags().Lookup("project"))

	gcpCmd.PersistentFlags().StringVarP(&credentials, "credentials", "c", "",
		"Path to GCP credentials JSON")
	gcpCmd.MarkPersistentFlagRequired("credentials")
	viper.BindPFlag("gcp_credentials", gcpCmd.PersistentFlags().Lookup("credentials"))
}

func setup() error {
	// Set provider config.
	viper.Set("provider", provider)

	// Setup the provider's image.
	providerVersion := "1.1.0"
	image := "gcr.io/kuda-project/provider-" + provider + ":" + providerVersion
	viper.Set("image", image)

	// Setup the volume mounting for the credentials.
	volumeSecret := docker.VolumeMapping{
		From: filepath.Dir(credentials),
		To:   "/secret",
	}
	viper.Set("volumes", []docker.VolumeMapping{volumeSecret})

	Setup()

	return nil
}
