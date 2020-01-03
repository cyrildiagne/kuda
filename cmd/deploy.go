package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/spf13/cobra"
)

// deployCmd represents the `kuda deploy` command.
var deployCmd = &cobra.Command{
	Use:   "deploy <manifest=./kuda.yaml>",
	Short: "Deploy the API remotely in production mode.",
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
		if err := deployWithSkaffold(manifest); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

func deployWithSkaffold(manifestFile string) error {
	// Load the manifest.
	manifest := latest.Manifest{}
	if err := loadManifest(manifestFile, &manifest); err != nil {
		return err
	}

	name := manifest.Name
	folder := cfg.Deployer.Skaffold.ConfigFolder

	skaffoldFile, err := generateSkaffoldConfigFiles(manifest.Deploy, name, folder)
	if err != nil {
		return err
	}

	// Run command.
	args := []string{"run", "-f", skaffoldFile}
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
