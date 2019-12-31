package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// deployCmd represents the `kuda deploy` command.
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys the API in production mode using Skaffold.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := deploy(); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

func deploy() error {
	args := []string{"run", "-f", "./.kuda/skaffold.yml"}

	// Run command.
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
