package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/api"
	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	skaffoldv1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
)

// generateCmd represents the `kuda deploy` command.
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the deployment files locally.",
	Run: func(cmd *cobra.Command, args []string) {
		registry, err := cmd.Flags().GetString("registry")
		if err != nil {
			panic(err)
		}

		folder, _ := cmd.Flags().GetString("to")
		if err := generate(folder, registry); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringP("registry", "r", "", "Registry where the image is pushed (Eg: gcr.io/kuda-project).")
	generateCmd.MarkFlagRequired("registry")
	generateCmd.Flags().StringP("to", "t", ".kuda", "Destination folder (default: .kuda)")
}

func generate(folder string, registry string) error {
	// Load the manifest.
	manifestFile := "./kuda.yaml"
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		fmt.Println("Could not load manifest", manifestFile)
		return err
	}

	// Make sure output folder exists.
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, 0700)
	}

	// Generate the files.
	im := api.ImageName{
		Author: cfg.Namespace,
		Name:   manifest.Name,
	}

	buildType := &skaffoldv1.BuildType{
		LocalBuild: &skaffoldv1.LocalBuild{},
	}

	// Generate Skaffold & Knative config files.
	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      cfg.Namespace,
		DockerArtifact: registry + "/" + im.GetID(),
		BuildType:      buildType,
	}

	// Add a release version of the manifest.
	manifest.Release = service.DockerArtifact
	manifYAML, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	manifYAML = []byte("# Generated automatically\n" + string(manifYAML))
	manifFile := filepath.FromSlash(folder + "/kuda.yaml")
	if err := utils.WriteYAML(manifYAML, manifFile); err != nil {
		return err
	}

	// Export API version in an env var for Skaffold's tagger.
	if err := utils.GenerateSkaffoldConfigFiles(service, manifest.Dev, folder+"/dev"); err != nil {
		return err
	}
	if err := utils.GenerateSkaffoldConfigFiles(service, manifest.Deploy, folder+"/deploy"); err != nil {
		return err
	}

	fmt.Printf("Files generated in %s\n", folder)

	return nil
}
