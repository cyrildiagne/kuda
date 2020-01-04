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

		if cfg.Deployer.Remote != nil {
			panic("dev is not yet supported on remote deployers")
		} else if cfg.Deployer.Skaffold != nil {
			// Start dev with Skaffold.
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if err := devWithSkaffold(*manifest, dryRun); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
	devCmd.Flags().Bool("dry-run", false, "Just generate the config files.")
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
