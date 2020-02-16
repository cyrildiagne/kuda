package deploy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cyrildiagne/kuda/pkg/api"
	"github.com/cyrildiagne/kuda/pkg/config"
	"github.com/cyrildiagne/kuda/pkg/utils"

	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
)

func deployFromReleaseManifest(manifest *latest.Manifest, env *api.Env, w http.ResponseWriter, r *http.Request) error {
	// Retrieve namespace.
	namespace, err := api.GetAuthorizedNamespace(env, r)
	if err != nil {
		return err
	}

	// Generate Knative YAML with appropriate namespace.
	service := config.ServiceSummary{
		Name:           manifest.Name,
		Namespace:      namespace,
		DockerArtifact: manifest.Release,
	}
	knativeCfg, err := config.GenerateKnativeConfig(service, manifest.Deploy)
	if err != nil {
		return err
	}
	knativeYAML, err := config.MarshalKnativeConfig(knativeCfg)
	if err != nil {
		return err
	}
	// Create new temp directory.
	tempDir, err := ioutil.TempDir("", namespace+"__"+manifest.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	knativeFile := filepath.FromSlash(tempDir + "/knative.yaml")
	if err := utils.WriteYAML(knativeYAML, knativeFile); err != nil {
		return err
	}

	// Run kubectl apply.
	args := []string{"apply", "-f", knativeFile}
	if err := RunCMD(w, "kubectl", args); err != nil {
		return err
	}

	// TODO: Add to the namespaces' deployments.

	return nil
}

func deployFromPublished(fromPublished string, env *api.Env, w http.ResponseWriter, r *http.Request) error {
	// Retrieve namespace.
	namespace, err := api.GetAuthorizedNamespace(env, r)
	if err != nil {
		return err
	}

	// Parse fromPublished to get author, name & version.
	im := api.ImageName{}
	if err := im.ParseFrom(fromPublished); err != nil {
		return err
	}

	// Check if image@version exists and is public.
	template, err := env.DB.GetVersionnedAPI(im)
	if err != nil {
		return err
	}
	if !template.IsPublic {
		err := fmt.Errorf("%s not found or not available", fromPublished)
		return api.StatusError{Code: 400, Err: err}
	}

	// Generate Knative YAML with appropriate namespace.
	service := config.ServiceSummary{
		Name:           im.Name,
		Namespace:      namespace,
		DockerArtifact: env.ContainerRegistry.GetDockerImagePath(im),
	}
	knativeCfg, err := config.GenerateKnativeConfig(service, template.Manifest.Deploy)
	if err != nil {
		return err
	}
	knativeYAML, err := config.MarshalKnativeConfig(knativeCfg)
	if err != nil {
		return err
	}
	// Create new temp directory.
	tempDir, err := ioutil.TempDir("", namespace+"__"+im.Name)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	knativeFile := filepath.FromSlash(tempDir + "/knative.yaml")
	if err := utils.WriteYAML(knativeYAML, knativeFile); err != nil {
		return err
	}

	// Run kubectl apply.
	args := []string{"apply", "-f", knativeFile}
	if err := RunCMD(w, "kubectl", args); err != nil {
		return err
	}

	// TODO: Add to the namespaces' deployments.

	return nil
}

func deployFromFiles(env *api.Env, w http.ResponseWriter, r *http.Request) error {
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

	// // Build with Skaffold.
	if err := Skaffold("run", contextDir, contextDir+"/skaffold.yaml", w); err != nil {
		return err
	}

	// Load the manifest.
	manifestFile := filepath.FromSlash(contextDir + "/kuda.yaml")
	manifest, err := utils.LoadManifest(manifestFile)
	if err != nil {
		return api.StatusError{Code: 400, Err: err}
	}

	// Register Template.
	apiVersion := &api.Version{
		IsPublic: false,
		Version:  manifest.Version,
		Manifest: manifest,
	}
	if err := registerAPI(env, namespace, apiVersion); err != nil {
		return err
	}

	// TODO: Add to the namespaces' deployments.

	fmt.Fprintf(w, "Deployment successful!\n")
	return nil
}

// HandleDeploy handles deployments from tar archived in body & published images.
func HandleDeploy(env *api.Env, w http.ResponseWriter, r *http.Request) error {
	// Set maximum upload size to 2GB.
	r.ParseMultipartForm((2 * 1000) << 20)

	fmt.Println("handle deploy")

	// Check if deploying from published.
	from := r.FormValue("from")
	if from != "" {
		return deployFromPublished(from, env, w, r)
	}

	// Otherwise check if deploying from a release manifest.
	release := r.FormValue("from-release")
	if release != "" {
		manifest := latest.Manifest{}
		if err := manifest.Load(strings.NewReader(release)); err != nil {
			return err
		}
		return deployFromReleaseManifest(&manifest, env, w, r)
	}

	// Otherwise try to deploy from the files attached.
	return deployFromFiles(env, w, r)
}
