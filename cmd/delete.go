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

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the `delete` command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete the remote clusters.",
	Long:  "Delete the remote clusters.",
	Run: func(cmd *cobra.Command, args []string) {
		Delete()
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}

// Delete a cluster.
func Delete() error {

	// Ask if we should delete the cluster.
	fmt.Print("Warning: This will delete the current cluster. Continue? (y/n) ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	txt := scanner.Text()
	if txt == "n" {
		return nil
	}

	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_delete"}
	// Run.
	err := RunDockerWithEnvs(docker.CommandOption{Image: image, Command: command})

	// Delete config file.
	config := viper.GetString("config")
	os.Remove(config)

	return err
}
