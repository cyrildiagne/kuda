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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the `get` command
var getCmd = &cobra.Command{
	Use:   "get [property]",
	Short: "Get information about the remote cluster.",
	Long:  "Get information about the remote cluster.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := get(args[0]); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func get(property string) error {
	fmt.Printf("→ Getting %s...\n", property)
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_get", property}

	// Run the command.
	dockerErr := RunDockerWithEnvs(docker.CommandOption{
		Image:   image,
		Command: command,
	})
	if dockerErr != nil {
		fmt.Println(dockerErr)
	}

	return nil
}
