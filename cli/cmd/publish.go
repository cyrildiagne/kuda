package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
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
		if err := publish(manifest); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
	publishCmd.Flags().Bool("dry-run", false, "Check the manifest for publication but skip execution.")
}

func publish(manifest *latest.Manifest) error {
	// Make sure a version is set.
	if manifest.Version == "" {
		return errors.New("version missing in manifest file")
	}
	// Create request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Add namespace
	writer.WriteField("namespace", cfg.Namespace)
	// Add manifest
	manifestYAML, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	writer.WriteField("manifest", string(manifestYAML))
	// Close writer
	writer.Close()

	url := cfg.Deployer.Remote.DeployerURL + "/publish"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err := sendToRemoteDeployer(req); err != nil {
		return err
	}

	return nil
}
