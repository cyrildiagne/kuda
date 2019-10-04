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
	"os"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appDeleteCmd represents the `app delete` command
var appDeleteCmd = &cobra.Command{
	Use:   "delete [app-name:app-version]",
	Short: "Deploy an app.",
	Long:  "Deploy an app.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := delete(args[0]); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	appCmd.AddCommand(appDeleteCmd)
}

func delete(app string) error {
	fmt.Println("→ Deleting app...")
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_app_delete", app}

	// Add the CWD to the volumes mounted in Docker.
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	volumes := []string{dir + ":/app_home"}

	// Run the command.
	dockerErr := RunDockerWithProviderEnvs(docker.CommandOption{
		Image:         image,
		Command:       command,
		AppendVolumes: volumes,
	})
	if dockerErr != nil {
		fmt.Println(dockerErr)
	}

	return nil
}
