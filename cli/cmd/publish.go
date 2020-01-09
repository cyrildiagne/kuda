package cmd

import (
	"fmt"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
)

// publishCmd represents the `kuda publish` command.
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Set an API image as publicly accessible. This doesn't affect your deployed APIs.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load the manifest
		manifestFile := "./kuda.yaml"
		manifest, err := utils.LoadManifest(manifestFile)
		if err != nil {
			fmt.Println("Could not load manifest", manifestFile)
			panic(err)
		}
		publish(manifest)
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
	publishCmd.Flags().Bool("dry-run", false, "Check the manifest for publication but skip execution.")
}

func publish(manifest *latest.Manifest) error {
	fmt.Println(manifest)
	// Check meta data in manifest:
	// Make sure Version is set
	return nil
}
