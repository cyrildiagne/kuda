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
	Use:   "dev",
	Short: "Deploy the API remotely in dev mode.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the manifest
		manifestFile := "./kuda.yaml"
		manifest, err := utils.LoadManifest(manifestFile)
		if err != nil {
			fmt.Println("Could not load manifest", manifestFile)
			panic(err)
		}

		// Check if dry run
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			panic(err)
		}

		if cfg.Deployer.Remote != nil {
			panic("dev is not yet supported on remote deployers")
		} else if cfg.Deployer.Skaffold != nil {
			// Start dev with Skaffold.
			if err := devWithSkaffold(*manifest, dryRun); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
	devCmd.Flags().Bool("dry-run", false, "Generate the config files but skip execution.")
}

func devWithSkaffold(manifest latest.Manifest, dryRun bool) error {

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
	fmt.Println("Config files have been written in:", folder)

	// Stop here if dry run.
	if dryRun {
		fmt.Println("Dry run: Skipping execution.")
		return nil
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
