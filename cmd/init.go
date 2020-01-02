package cmd

import (
	"fmt"
	"os"

	"github.com/cyrildiagne/kuda/pkg/kuda/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// initCmd represents the `kuda init` command.
var initCmd = &cobra.Command{
	Use:   "init <deployer>",
	Short: "Initializes the local configuration.",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		deployer := args[0]

		// Handle skaffold provider.
		if deployer == "skaffold" {

			// Ensure that users provide the docker_registry when using the skaffold
			// deployer.
			dockerRegistry, e := cmd.Flags().GetString("docker_registry")
			if dockerRegistry == "" || e != nil {
				panic("The skaffold deployer requires a [-d, --docker_registry] value.")
			}

			// Create a Kuda config.
			var newCfg config.UserConfig
			newCfg.Namespace, e = cmd.Flags().GetString("namespace")
			if e != nil {
				panic("Could not retrieve namespace flag.")
			}

			// Setup the skaffold config.
			if deployer == "skaffold" {
				newCfg.Deployer.Skaffold = &config.SkaffoldDeployerConfig{
					DockerRegistry: dockerRegistry,
					ConfigFolder:   "./.kuda",
				}
			}

			// Write the file to disk.
			writeConfig(newCfg)

		} else {
			panic("Only 'skaffold' is currently supported as deployer.")
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("namespace", "n", "default", "Knative namespace.")
	initCmd.Flags().StringP("docker_registry", "d", "", "Docker registry.")
}

func writeConfig(cfg config.UserConfig) error {
	content, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	f, err := os.Create(cfgFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(content); err != nil {
		return err
	}
	fmt.Println("Config written in " + cfgFile)
	return nil
}
