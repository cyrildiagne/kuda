package deploy

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cyrildiagne/kuda/pkg/api"
	"github.com/cyrildiagne/kuda/pkg/utils"
)

func registerAPI(env *api.Env, author string, template *api.Version) error {
	im := api.ImageName{
		Author:  author,
		Name:    template.Manifest.Name,
		Version: template.Version,
	}

	// Update API Metadata.
	metadata := &map[string]interface{}{
		"author": im.Author,
		"name":   im.Name,
		"image":  env.ContainerRegistry.GetDockerImagePath(im),
	}
	if err := env.DB.UpdateAPIMetadata(im.GetID(), metadata); err != nil {
		return err
	}

	// Update api version document
	if err := env.DB.UpdateVersionnedAPI(im.GetID(), im.Version, template); err != nil {
		return err
	}

	return nil
}

// HandlePublish publishes from tar file in body.
func HandlePublish(env *api.Env, w http.ResponseWriter, r *http.Request) error {
	// Set maximum upload size to 2GB.
	r.ParseMultipartForm((2 * 1000) << 20)

	// Retrieve namespace.
	namespace, err := api.GetAuthorizedNamespace(env, r)
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
		return api.StatusError{Code: 400, Err: err}
	}

	// Register API.
	apiVersion := &api.Version{
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
