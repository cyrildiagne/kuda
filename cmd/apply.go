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
	"errors"
	"fmt"
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/docker"
	"github.com/spf13/cobra"
)

// applyCmd represents the `apply` command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a config file to the remote clusters.",
	Long:  "Apply a config file to the remote clusters.",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := cmd.Flags().GetString("filename")
		if err != nil {
			panic(errors.New("cannot retrieve --filename | -f parameter"))
		}
		Apply(file)
	},
}

func init() {
	RootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringP("filename", "f", "", "Config file to apply.")
	applyCmd.MarkFlagRequired("filename")
}

// Apply a cluster.
func Apply(file string) error {
	// Mount the folder containing the config file in Docker.
	filename := filepath.Base(file)
	filedir, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		fmt.Println("Config file path is invalid.")
		panic(err)
	}
	volumes := []string{
		filedir + ":/config",
	}
	// Command to run.
	command := []string{"kuda_apply", "/config/" + filename}
	// Run.
	dockerErr := RunProviderCommand(docker.CommandOption{
		Command:       command,
		AppendVolumes: volumes,
	})

	return dockerErr
}
