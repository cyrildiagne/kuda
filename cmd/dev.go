package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/spf13/cobra"
)

// devCmd represents the `kuda dev` command.
var devCmd = &cobra.Command{
	Use:   "dev <manifest=./kuda.yaml>",
	Short: "Deploy the API remotely in dev mode.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manifest := "./kuda.yaml"
		if len(args) == 1 {
			manifest = args[0]
		}
		// Ensure manifest exists
		if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
			fmt.Println("Could not load manifest", manifest)
			panic(err)
		}
		if err := devWithSkaffold(manifest); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}

func devWithSkaffold(manifestFile string) error {
	// Load the manifest.
	manifest := latest.Manifest{}
	if err := loadManifest(manifestFile, &manifest); err != nil {
		return err
	}

	name := manifest.Name + "-dev"
	folder := cfg.Deployer.Skaffold.ConfigFolder

	skaffoldFile, err := generateSkaffoldConfigFiles(manifest.Dev, name, folder)
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
