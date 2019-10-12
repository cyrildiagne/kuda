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
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setupCmd represents the `setup` command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup a remote clusters",
	Long:  "Setup a remote clusters.",
}

func init() {
	RootCmd.AddCommand(setupCmd)
}

// Setup is a unified setup function across all providers.
func Setup() error {
	// Image to run.
	image := viper.GetString("image")
	// Command to run.
	command := []string{"kuda_setup"}
	// Run
	err := RunDockerWithEnvs(docker.CommandOption{Image: image, Command: command})
	if err != nil {
		// Ask if we should delete the cluster.
		fmt.Print("There was an error setting up the cluster. Do you want to delete it? (y/n) ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		txt := scanner.Text()
		if txt == "y" {
			Delete()
		}
	}

	// Write new config to home directory.
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	cfgFile := filepath.FromSlash(home + "/.kuda.yaml")
	viper.SetConfigFile(cfgFile)
	viper.WriteConfig()
	fmt.Println("Config written in " + viper.ConfigFileUsed())

	return nil
}
