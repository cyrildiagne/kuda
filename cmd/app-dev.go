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
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// appDevCmd represents the `app dev` command
var appDevCmd = &cobra.Command{
	Use:   "dev [app-name] [app-dir]",
	Short: "Deploy an app in dev mode.",
	Long:  "Deploy an app in dev mode.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Set current working directory from 2nd argument if provided otherwise
		// use the current working directory.
		dir := ""
		if len(args) > 1 {
			argDir, err := filepath.Abs(args[1])
			if err != nil {
				panic(err)
			}
			dir = argDir
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			dir = cwd
		}

		if err := dev(args[0], dir); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	appCmd.AddCommand(appDevCmd)
}

func dev(appName string, appDir string) error {
	fmt.Println("→ Start app dev...")
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_app_dev", appName}
	// Add the application folder to the volumes mounted in Docker.
	volumes := []string{
		// Bind the app home directory.
		appDir + ":/app_home",
		// Bind docker socker for Skaffold.
		"/var/run/docker.sock:/var/run/docker.sock",
	}
	// Run the command.
	dockerErr := RunDockerWithEnvs(docker.CommandOption{
		Image:         image,
		Command:       command,
		AppendVolumes: volumes,
	})
	if dockerErr != nil {
		fmt.Println(dockerErr)
	}

	return nil
}
