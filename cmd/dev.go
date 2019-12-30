package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// devCmd represents the `setup init` command.
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs the API remotely in dev mode using Skaffold.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := dev(); err != nil {
			fmt.Println("ERROR:", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}

func dev() error {
	args := []string{"dev", "-f", "./.kuda/skaffold-dev.yml"}

	// Run command.
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	return err
}
