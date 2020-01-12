package main

import (
	"fmt"
	"os"

	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// initCmd represents the `kuda init` command.
var initCmd = &cobra.Command{
	Use:   "init <namespace>",
	Short: "Initializes the local configuration.",
	Args:  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Create a Kuda config.
		var newCfg config.UserConfig
		newCfg.Namespace = args[0]

		// Retrieve provider
		provider, err := cmd.Flags().GetString("provider")
		if err != nil {
			panic(err)
		}
		// Setup the default remote config.
		authURL := "https://auth." + provider
		authURLFlag, _ := cmd.Flags().GetString("auth_url")
		if authURLFlag != "" {
			authURL = authURLFlag
		}
		apiURL := "https://api." + provider
		apiURLFlag, _ := cmd.Flags().GetString("api_url")
		if apiURLFlag != "" {
			apiURL = apiURLFlag
		}
		newCfg.Provider = config.ProviderConfig{
			AuthURL: authURL,
			ApiURL:  apiURL,
		}

		// Start login flow.
		fmt.Println("Authenticating on...", newCfg.Provider.AuthURL)
		user, err := startLoginFlow(newCfg.Provider.AuthURL)
		if err != nil {
			fmt.Println("Authentication error.")
			panic(err)
		}
		newCfg.Provider.User = user
		fmt.Println("Authenticated as", user.DisplayName)

		// Write the file to disk.
		writeConfig(newCfg)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("provider", "p", "kuda.cloud", "Knative namespace.")
	initCmd.Flags().String("auth_url", "", "Authentication URL.")
	initCmd.Flags().String("api_url", "", "Deployer URL.")
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
