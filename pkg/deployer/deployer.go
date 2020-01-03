package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"github.com/gorilla/mux"
)

var gcpProjectID string
var dockerRegistry string

func handleDeployment(w http.ResponseWriter, r *http.Request) {
	// Set maximum upload size to 2GB.
	r.ParseMultipartForm((2 * 1000) << 20)

	// Retrieve namespace.
	namespace := r.FormValue("namespace")
	if namespace == "" {
		http.Error(w, "error retrieving namespace", 500)
		return
	}

	// Retrieve Filename, Header and Size of the file.
	file, handler, err := r.FormFile("context")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error retrieving file", 500)
		return
	}
	defer file.Close()
	fmt.Printf("File: %+v, %+v Ko\n", handler.Filename, handler.Size/1024)
	fmt.Printf("Header: %+v\n", handler.Header)

	// Create new temp directory.
	tempDir, err := ioutil.TempDir("", namespace)
	fmt.Println(tempDir)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error creating temp dir", 500)
		return
	}
	defer os.RemoveAll(tempDir) // Clean up.

	// Extract file to temp directory.
	err = utils.Untar(tempDir, file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error extracting content", 500)
		return
	}

	// Load the manifest.
	manifestFile := filepath.FromSlash(tempDir + "/kuda.yaml")
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not load manifest", 500)
		return
	}

	dockerArtifact := dockerRegistry + "/" + namespace + "__" + manifest.Name

	// Generate Skaffold & Knative config files.
	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      namespace,
		DockerArtifact: dockerArtifact,
		BuildType: v1.BuildType{
			GoogleCloudBuild: &v1.GoogleCloudBuild{
				ProjectID: gcpProjectID,
			},
		},
	}
	folder := filepath.FromSlash(tempDir + "/.kuda")
	skaffoldFile, err := utils.GenerateSkaffoldConfigFiles(service, manifest.Deploy, folder)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not generate config files", 500)
		return
	}

	// Run Skaffold Deploy.
	args := []string{"run", "-f", skaffoldFile}
	cmd := exec.Command("skaffold", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Dir = tempDir

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		http.Error(w, "error running skaffold", 500)
		return
	}

	fmt.Fprintf(w, "Deployment successful!\n")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello!\n")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initKubectlGcloud() {
	// Authenticate gcloud using application credentials.
	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file",
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("Error authenticating with credentials.")
	}

	// Get kubeconfig.
	args := []string{"container", "clusters", "get-credentials",
		"--project", gcpProjectID,
		"--region", "us-central1-a", "kuda"}
	cmd = exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("could not retrieve kubectl credentials")
	}
}

func main() {
	gcpProjectID = os.Getenv("KUDA_GCP_PROJECT")
	if gcpProjectID == "" {
		panic("cloud not load env var KUDA_GCP_PROJECT")
	}
	log.Println("Using project:", gcpProjectID)

	dockerRegistry = "gcr.io/" + gcpProjectID
	log.Println("Using registry:", dockerRegistry)

	initKubectlGcloud()

	port := getEnv("PORT", "8080")
	fmt.Println("Starting deployer on port", port)

	r := mux.NewRouter()
	r.HandleFunc("/", hello).Methods("GET")
	r.HandleFunc("/", handleDeployment).Methods("POST")
	http.ListenAndServe(":"+port, r)
}
