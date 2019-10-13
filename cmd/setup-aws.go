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

	"github.com/cyrildiagne/kuda/pkg/docker"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// awsCmd represents the `setup gcp` command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "Setup Kuda on AWS.",
	Long:  "Setup Kuda on AWS.",
	Run: func(cmd *cobra.Command, args []string) {
		viper.BindPFlags(cmd.Flags())
		if err := setupAWS(); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	setupCmd.AddCommand(awsCmd)

	// viper.BindPFlags(gcpCmd.Flags())
}

func setupAWS() error {
	const provider = "aws"
	const providerVersion = "0.1.0"

	// Set provider config.
	viper.Set("provider", provider)

	// Setup the provider's image.
	image := "gcr.io/kuda-project/provider-" + provider + ":" + providerVersion
	viper.Set("image", image)

	// Setup the volume mounting for the credentials.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	volumeSecret := docker.VolumeMapping{
		From: home + "/.aws",
		To:   "/aws-credentials/",
	}
	viper.Set("volumes", []docker.VolumeMapping{volumeSecret})

	Setup()

	return nil
}
