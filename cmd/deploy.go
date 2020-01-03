package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"

	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/spf13/cobra"
)

// deployCmd represents the `kuda deploy` command.
var deployCmd = &cobra.Command{
	Use:   "deploy <manifest=./kuda.yaml>",
	Short: "Deploy the API remotely in production mode.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		manifestFile := "./kuda.yaml"
		if len(args) == 1 {
			manifestFile = args[0]
		}
		// Load the manifest
		manifest, err := utils.LoadManifest(manifestFile)
		if err != nil {
			fmt.Println("Could not load manifest", manifestFile)
			panic(err)
		}

		if cfg.Deployer.Remote != nil {
			if err := deployWithRemote(manifest); err != nil {
				fmt.Println("ERROR:", err)
			}
		} else if cfg.Deployer.Skaffold != nil {
			if err := deployWithSkaffold(manifest); err != nil {
				fmt.Println("ERROR:", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(deployCmd)
}

func deployWithRemote(manifest *latest.Manifest) error {
	// Create destination tar file
	output, err := ioutil.TempFile("", "*.rar")
	fmt.Println("Building context tar:", output.Name())
	if err != nil {
		return err
	}
	defer os.Remove(output.Name())

	// Open .dockerignore file if it exists
	dockerignore, err := os.Open(".dockerignore")
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer dockerignore.Close()

	// Tar context folder.
	utils.Tar("./", output.Name(), output, dockerignore)

	fmt.Println("Sending to deployer:", cfg.Deployer.Remote.URL)

	// Create request
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

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

	// Add namespace
	writer.WriteField("namespace", cfg.Namespace)

	// Close writer
	writer.Close()

	url := cfg.Deployer.Remote.URL
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Check response.
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		fmt.Println(resp.Status)
		fmt.Println(string(respBody))
		return fmt.Errorf("error with remote deployer")
	}
	fmt.Println(string(respBody))

	return nil
}

func deployWithSkaffold(manifest *latest.Manifest) error {

	folder := cfg.Deployer.Skaffold.ConfigFolder
	registry := cfg.Deployer.Skaffold.DockerRegistry

	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      cfg.Namespace,
		DockerArtifact: registry + "/" + manifest.Name,
	}

	skaffoldFile, err := utils.GenerateSkaffoldConfigFiles(service, manifest.Deploy, folder)
	if err != nil {
		return err
	}

	// Run command.
	args := []string{"run", "-f", skaffoldFile}
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
