package cmd

import (
	"fmt"
	"os"

	"github.com/cyrildiagne/kuda/pkg/kuda"

	"github.com/spf13/cobra"
)

// initCmd represents the `kuda init` command.
var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Generate the configuration files in a local .kuda folder.",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		dockerRegistry, _ := cmd.Flags().GetString("docker_registry")
		namespace, _ := cmd.Flags().GetString("namespace")

		err := checkFolder()
		if err != nil {
			panic(err)
		}

		// Create config.
		cfg := kuda.NewConfig(name, namespace)
		cfg.DockerDestImage = dockerRegistry
		err = generate(cfg)
		if err != nil {
			panic(err)
		}

		// Create dev config
		cfgDev := kuda.NewConfig(name, namespace)
		cfgDev.DockerDestImage = dockerRegistry
		cfgDev.AddDevConfigFlask()
		cfgDev.SetFilesSuffix("-dev")
		err = generate(cfgDev)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("docker_registry", "d", "", "Docker registry.")
	initCmd.MarkFlagRequired("docker_registry")

	initCmd.Flags().StringP("namespace", "n", "default", "Knative namespace.")
}

func checkFolder() error {
	if _, err := os.Stat("./main.py"); os.IsNotExist(err) {
		fmt.Println("WARNING: Folder does not contain a main.py file." +
			" Edit `.kuda/service-dev.yml` to enable `kuda dev`")
	}

	// fmt.Println(files)
	return nil
}

func generate(cfg kuda.Config) error {
	manifests, err := kuda.GenerateConfigFiles(cfg)
	if err != nil {
		return err
	}

	// Make sure config folders exists.
	if _, err := os.Stat(cfg.ConfigFolder); os.IsNotExist(err) {
		os.Mkdir(cfg.ConfigFolder, 0700)
	}

	// Write dev manifests.
	writeManifest(manifests.Kservice, cfg.KserviceFile)
	writeManifest(manifests.Skaffold, cfg.SkaffoldFile)

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
