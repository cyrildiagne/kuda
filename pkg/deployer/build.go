package deployer

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/utils"
)

func generate(namespace string, contextDir string, env *Env) error {
	// Load the manifest.
	manifestFile := filepath.FromSlash(contextDir + "/kuda.yaml")
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return StatusError{400, err}
	}

	// TODO: replace namespace by user ID.
	dockerArtifact := env.GetDockerImagePath(namespace, manifest.Name)

	// Generate Skaffold & Knative config files.
	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      namespace,
		DockerArtifact: dockerArtifact,
		BuildType: v1.BuildType{
			GoogleCloudBuild: &v1.GoogleCloudBuild{
				ProjectID: env.GCPProjectID,
			},
		},
	}
	// Export API version in an env var for Skaffold's tagger.
	os.Setenv("API_VERSION", manifest.Version)
	if err := utils.GenerateSkaffoldConfigFiles(service, manifest.Deploy, contextDir); err != nil {
		return err
	}
	return nil
}

func extractContext(prefix string, r *http.Request) (string, error) {
	// Retrieve Filename, Header and Size of the file.
	file, _, err := r.FormFile("context")
	if err != nil {
		return "", err
	}
	defer file.Close()
	// Create new temp directory.
	tempDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		return "", err
	}
	// Extract file to temp directory.
	err = utils.Untar(tempDir, file)
	if err != nil {
		return "", err
	}
	// Return tempDir path
	return tempDir, nil
}
