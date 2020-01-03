package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
)

// devCmd represents the `kuda dev` command.
var devCmd = &cobra.Command{
	Use:   "dev <manifest=./kuda.yaml>",
	Short: "Deploy the API remotely in dev mode.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manifestFile := "./kuda.yaml"
		if len(args) == 1 {
			manifestFile = args[0]
		}
		// Load the manifest
		manifest, err := utils.LoadManifest(manifestFile)
		if err != nil {
			fmt.Println("Could not load manifest", manifestFile)
			panic(err)
		}
		// Start dev with Skaffold.
		if err := devWithSkaffold(*manifest); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}

func devWithSkaffold(manifest latest.Manifest) error {

	folder := cfg.Deployer.Skaffold.ConfigFolder
	registry := cfg.Deployer.Skaffold.DockerRegistry

	service := config.ServiceSummary{
		Name:           manifest.Name + "-dev",
		Namespace:      cfg.Namespace,
		DockerArtifact: registry + "/" + manifest.Name,
	}

	skaffoldFile, err := utils.GenerateSkaffoldConfigFiles(service, manifest.Dev, folder)
	if err != nil {
		return err
	}

	// Run command.
	args := []string{"dev", "-f", skaffoldFile}
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
