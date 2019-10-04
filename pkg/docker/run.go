/*
Package docker -

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
package docker

import (
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

// CommandOption Options to run the RunDockerCommand.
type CommandOption struct {
	Image   string
	Command []string
	// Extra values to be appended to viper's config.
	AppendVolumes []string
	AppendEnv     []string
}

// VolumeMapping maps local volumes to mount in Docker.
// type Volume map[string]string
type VolumeMapping struct {
	From string
	To   string
}

// RunDockerCommand runs a docker command
// using envs & volumes from viper's config.
func RunDockerCommand(opts CommandOption) error {
	// Run docker commands in interactive mode and with TTY with -it.
	// Also remove the container as soon as it's been ran with --rm.
	args := []string{"run", "-it", "--rm"}

	// Environment Variables.
	for _, e := range opts.AppendEnv {
		args = append(args, "-e", e)
	}

	// Volume mappings.
	var configVols []VolumeMapping
	viper.UnmarshalKey("volumes", &configVols)
	for _, v := range configVols {
		args = append(args, "-v", v.From+":"+v.To)
	}
	for _, v := range opts.AppendVolumes {
		args = append(args, "-v", v)
	}

	// Set image & command.
	args = append(args, opts.Image)
	args = append(args, opts.Command...)

	// Run command.
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
