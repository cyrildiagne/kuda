package deployer

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/utils"
)

// HandlePublish publishes from tar file in body.
func HandlePublish(env *Env, w http.ResponseWriter, r *http.Request) error {
	// Set maximum upload size to 2GB.
	r.ParseMultipartForm((2 * 1000) << 20)

	// Retrieve namespace.
	namespace, err := GetAuthorizedNamespace(env, r)
	if err != nil {
		return err
	}

	// Extract archive to temp folder.
	contextDir, err := extractContext(namespace, r)
	if err != nil {
		return err
	}
	defer os.RemoveAll(contextDir) // Clean up.

	// Build and push image.
	if err := generate(namespace, contextDir, env); err != nil {
		return err
	}

	// Setup client stream.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/event-stream")

	// Build with Skaffold.
	if err := Skaffold("build", contextDir, contextDir+"/skaffold.yaml", w); err != nil {
		return err
	}

	// Load the manifest.
	manifestFile := filepath.FromSlash(contextDir + "/kuda.yaml")
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return StatusError{400, err}
	}

	// Register API.
	apiVersion := APIVersion{
		IsPublic: true,
		Version:  manifest.Version,
		Manifest: manifest,
	}
	if err := registerAPI(env, namespace, apiVersion); err != nil {
		return err
	}

	fmt.Fprintf(w, "Publish successful!\n")
	return nil
}
