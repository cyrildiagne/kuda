package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	skaffoldCfg "github.com/cyrildiagne/kuda/pkg/kuda/deployer/skaffold/config"
	"github.com/cyrildiagne/kuda/pkg/kuda/manifest/latest"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
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

	cfgFolder := cfg.Deployer.Skaffold.ConfigFolder

	// Make sure output folder exists.
	if _, err := os.Stat(cfgFolder); os.IsNotExist(err) {
		os.Mkdir(cfgFolder, 0700)
	}

	// Generate the knative yaml file.
	knativeCfg, err := skaffoldCfg.GenerateKnativeConfig(manifest.Name, manifest.Deploy, cfg)
	if err != nil {
		return err
	}
	knativeYAML, err := skaffoldCfg.MarshalKnativeConfig(knativeCfg)
	if err != nil {
		return err
	}
	knativeFile := filepath.FromSlash(cfgFolder + "/knative.yaml")
	if err := writeYAML(knativeYAML, knativeFile); err != nil {
		return err
	}

	// Generate the skaffold yaml file.
	skaffoldCfg, err := skaffoldCfg.GenerateSkaffoldConfig(manifest.Name, manifest.Deploy, cfg, knativeFile)
	if err != nil {
		return err
	}
	skaffoldYAML, err := yaml.Marshal(skaffoldCfg)
	if err != nil {
		return err
	}
	skaffoldFile := filepath.FromSlash(cfgFolder + "/skaffold.yaml")
	if err := writeYAML(skaffoldYAML, skaffoldFile); err != nil {
		return err
	}

	// // Run command.
	// args := []string{"run", "-f", skaffoldFile}
	// cmd := exec.Command("skaffold", args...)
	// cmd.Stdout = os.Stdout
	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	return err
	// }

	return nil
}
