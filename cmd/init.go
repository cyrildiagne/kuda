package cmd

import (
	"fmt"
	"os"

	"github.com/cyrildiagne/kuda/pkg/config"
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

		// Create a Kuda config.
		var newCfg config.UserConfig
		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			panic("Could not retrieve the namespace flag.")
		}
		newCfg.Namespace = namespace

		// Handle skaffold deployer.
		if deployer == "skaffold" {

			// Ensure that users provide the docker_registry when using the skaffold
			// deployer.
			dockerRegistry, e := cmd.Flags().GetString("docker_registry")
			if dockerRegistry == "" || e != nil {
				panic("The skaffold deployer requires a [-d, --docker_registry] value.")
			}

			// Setup the skaffold config.
			newCfg.Deployer.Skaffold = &config.SkaffoldDeployerConfig{
				DockerRegistry: dockerRegistry,
				ConfigFolder:   "./.kuda",
			}

			// Write the file to disk.
			writeConfig(newCfg)

		} else {

			// Setup the default remote config.
			authURL := "https://auth.kuda." + deployer
			authURLFlag, _ := cmd.Flags().GetString("auth_url")
			if authURLFlag != "" {
				authURL = authURLFlag
			}
			deployerURL := "https://deployer.kuda." + deployer
			deployerURLFlag, _ := cmd.Flags().GetString("deployer_url")
			if deployerURLFlag != "" {
				deployerURL = deployerURLFlag
			}
			newCfg.Deployer.Remote = &config.RemoteDeployerConfig{
				AuthURL:     authURL,
				DeployerURL: deployerURL,
			}

			// Start login flow.
			fmt.Println("Authenticating on...", newCfg.Deployer.Remote.AuthURL)
			user, err := startLoginFlow(newCfg.Deployer.Remote.AuthURL)
			if err != nil {
				fmt.Println("Authentication error.")
				panic(err)
			}
			newCfg.Deployer.Remote.User = user
			fmt.Println("Authenticated as", user.DisplayName)

			// Write the file to disk.
			writeConfig(newCfg)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("namespace", "n", "default", "Knative namespace.")
	initCmd.Flags().StringP("docker_registry", "d", "", "Docker registry.")
	initCmd.Flags().String("auth_url", "", "Authentication URL.")
	initCmd.Flags().String("deployer_url", "", "Deployer URL.")
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
