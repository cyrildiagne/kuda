package cmd

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
	// "github.com/go-openapi/loads/fmts"
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
			fmt.Println("Could not load ./kuda.yaml", manifestFile)
			panic(err)
		}

		// TODO: Ensure there is an OpenAPI spec.

		// Publish
		if err := publish(manifest); err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(publishCmd)
}

func publish(manifest *latest.Manifest) error {
	// Create request body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Add context
	if err := addContextFilesToRequest("./", writer); err != nil {
		return err
	}
	// Add namespace
	writer.WriteField("namespace", cfg.Namespace)
	// Close writer
	writer.Close()

	// Create request.
	url := cfg.Deployer.Remote.DeployerURL + "/publish"
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send to remote deployer.
	if err := sendToRemoteDeployer(req); err != nil {
		return err
	}

	return nil
}
