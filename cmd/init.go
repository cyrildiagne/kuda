package cmd

import (
	"fmt"

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

		if err := kuda.GenerateConfigFiles(url, dockerRegistry); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("docker_registry", "d", "", "Docker registry.")
	initCmd.MarkFlagRequired("docker_registry")
}
