package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the `kuda deploy` command.
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the API remotely in production mode.",
	Run: func(cmd *cobra.Command, args []string) {
		published, _ := cmd.Flags().GetString("from")
		if published != "" {
			deployFromPublished(published)
		} else {
			deployFromLocal()
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringP("from", "f", "", "Fully qualified name of a published API image.")
}

func deployFromPublished(published string) error {
	fmt.Println("Deploy from published API image", published)

	fmt.Println("Sending to deployer:", cfg.Provider.DeployerURL)

	// Create request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// Add namespace
	writer.WriteField("namespace", cfg.Namespace)
	writer.WriteField("from_published", published)
	// Close writer
	writer.Close()

	url := cfg.Provider.DeployerURL + "/deploy"
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

func deployFromLocal() {
	// Load the manifest
	manifestFile := "./kuda.yaml"
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		fmt.Println("Could not load manifest", manifestFile)
		panic(err)
	}

	if err := deploy(manifest); err != nil {
		fmt.Println("ERROR:", err)
	}
}

func addContextFilesToRequest(source string, writer *multipart.Writer) error {
	// Create destination tar file
	output, err := ioutil.TempFile("", "*.tar")
	fmt.Println("Building context tar:", output.Name())
	if err != nil {
		return err
	}

	// Open .dockerignore file if it exists
	dockerignore, err := os.Open(".dockerignore")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer dockerignore.Close()

	// Tar context folder.
	utils.Tar(source, output.Name(), output, dockerignore)

	// Defer the deletion of the temp tar file.
	defer os.Remove(output.Name())

	// Add tar file to request
	file, err := os.Open(output.Name())
	defer file.Close()
	if err != nil {
		return err
	}
	part, err := writer.CreateFormFile("context", "context.tar")
	if err != nil {
		return err
	}
	io.Copy(part, file)

	return nil
}

func deploy(manifest *latest.Manifest) error {
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
	url := cfg.Provider.DeployerURL + "/deploy"
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

func sendToRemoteDeployer(req *http.Request) error {
	accessToken := "Bearer " + cfg.Provider.User.Token.AccessToken
	req.Header.Set("Authorization", accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read body stream.
	br := bufio.NewReader(resp.Body)
	for {
		bs, err := br.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Print(string(bs))
	}

	// Check response.
	if resp.StatusCode != 200 {
		fmt.Println("Sending to deployer returned an error", resp.Status)
		if resp.StatusCode == 401 {
			fmt.Println("Try authenticating again running 'kuda init <args>'.")
		}
		return fmt.Errorf("error with remote deployer")
	}
	return nil
}
