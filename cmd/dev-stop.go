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
	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// stopCmd represents the `dev stop` command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a dev session.",
	Long:  "Stop a dev session.",
	Run: func(cmd *cobra.Command, args []string) {
		Stop()
	},
}

func init() {
	devCmd.AddCommand(stopCmd)
}

// Stop the remote session.
func Stop() error {
	// color.Cyan("→ Stopping the remote session...")
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_dev_stop"}
	// Run
	err := RunDockerWithEnvs(docker.CommandOption{Image: image, Command: command})
	return err
}
