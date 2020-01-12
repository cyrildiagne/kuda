package main

import (
	"fmt"

	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
)

// devCmd represents the `kuda dev` command.
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Deploy the API remotely in dev mode.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the manifest
		manifestFile := "./kuda.yaml"
		manifest, err := utils.LoadManifest(manifestFile)
		if err != nil {
			fmt.Println("Could not load manifest", manifestFile)
			panic(err)
		}
		fmt.Println(manifest)

		panic("dev is not yet supported.")
	},
}

func init() {
	RootCmd.AddCommand(devCmd)
}
