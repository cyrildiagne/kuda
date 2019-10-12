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
var providerVersion = "1.2.0"

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

	gcpCmd.PersistentFlags().StringP("gcp_project_id", "p", "", "GCP Project ID")
	gcpCmd.MarkPersistentFlagRequired("gcp_project_id")

	gcpCmd.PersistentFlags().StringP("gcp_credentials", "c", "", "Path to GCP credentials JSON")
	gcpCmd.MarkPersistentFlagRequired("gcp_credentials")
	viper.BindPFlags(gcpCmd.PersistentFlags())

	gcpCmd.Flags().String("gcp_cluster_name", "kuda", "Name of the cluster.")
	gcpCmd.Flags().String("gcp_compute_zone", "us-central1-a", "Compute Zone for the cluster.")
	gcpCmd.Flags().String("gcp_machine_type", "n1-standard-4", "Machine type.")
	gcpCmd.Flags().Int("gcp_pool_num_nodes", 1, "Default number of nodes on the system pool. ")
	gcpCmd.Flags().String("gcp_gpu", "k80", "Default GPU to use")
	gcpCmd.Flags().Bool("gcp_use_preemptible", false, "Wether or not to use pre-emptible instances")
	viper.BindPFlags(gcpCmd.Flags())
}

func setup() error {
	// Set provider config.
	viper.Set("provider", provider)

	// Setup the provider's image.
	image := "gcr.io/kuda-project/provider-" + provider + ":" + providerVersion
	viper.Set("image", image)

	// Setup the volume mounting for the credentials.
	credentials := viper.GetString("gcp_credentials")
	volumeSecret := docker.VolumeMapping{
		From: filepath.Dir(credentials),
		To:   "/secret",
	}
	viper.Set("volumes", []docker.VolumeMapping{volumeSecret})

	Setup()

	return nil
}
