package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cyrildiagne/kuda/pkg/kuda"

	"github.com/spf13/cobra"
)

// initCmd represents the `kuda init` command.
var initCmd = &cobra.Command{
	Use:   "init <URL>",
	Short: "Generate the configuration files in a local .kuda folder.",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		dockerRegistry, _ := cmd.Flags().GetString("docker_registry")

		err := checkFolder()
		if err != nil {
			log.Fatal("ERROR:", err)
		}

		manifests, err := kuda.GenerateConfigFiles(url, dockerRegistry)
		if err != nil {
			log.Fatal("ERROR:", err)
		}

		// Write prod manifests.
		writeManifest(manifests.Prod.Kservice, manifests.Prod.Config.KserviceFile)
		writeManifest(manifests.Prod.Skaffold, manifests.Prod.Config.SkaffoldFile)
		// Write dev manifests.
		writeManifest(manifests.Dev.Kservice, manifests.Dev.Config.KserviceFile)
		writeManifest(manifests.Dev.Skaffold, manifests.Dev.Config.SkaffoldFile)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("docker_registry", "d", "", "Docker registry.")
	initCmd.MarkFlagRequired("docker_registry")
}

func checkFolder() error {
	if _, err := os.Stat("./main.py"); os.IsNotExist(err) {
		fmt.Println("WARNING: Folder does not contain a main.py file." +
			" Edit `.kuda/service-dev.yml` to enable `kuda dev`")
	}
	// fmt.Println(files)
	return nil
}

// WriteManifest writes the manifest to disk
func writeManifest(content string, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}
