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
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the `dev start` command
var startCmd = &cobra.Command{
	Use:   "start [docker-image] [dir]",
	Short: "Start a dev session.",
	Long:  "Start a dev session using the provider docker image.",
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

		if err := start(args[0], dir); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	devCmd.AddCommand(startCmd)
}

func start(devImage string, folderMount string) error {
	fmt.Println("→ Starting a remote session...")
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_dev_start", devImage}

	// Add the CWD to the volumes mounted in Docker.
	volumes := []string{folderMount + ":/app_home"}

	// Run the command.
	dockerErr := RunDockerWithEnvs(docker.CommandOption{
		Image:         image,
		Command:       command,
		AppendVolumes: volumes,
	})
	if dockerErr != nil {
		fmt.Println(dockerErr)
	}

	// Ask if we should stop the session.
	fmt.Print("Do you want to delete the remote session? (y/n) ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	txt := scanner.Text()
	if txt == "y" {
		Stop()
	}

	return nil
}
