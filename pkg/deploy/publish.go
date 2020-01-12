package deploy

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"cloud.google.com/go/firestore"
	"github.com/cyrildiagne/kuda/pkg/api"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
	"github.com/cyrildiagne/kuda/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIVersion stores an API version.
type APIVersion struct {
	IsPublic bool             `firestore:"isPublic"`
	Version  string           `firestore:"version"`
	Manifest *latest.Manifest `firestore:"manifest"`
	// Paths    openapi.Paths    `firestore:"paths,omitempty"`
	// Paths    *openapi3.Swagger   `firestore:"openapi,omitempty"`
	// Paths openapi3.Paths `firestore:"openapi,omitempty"`
	// Paths    map[string]*openapi3.PathItem `firestore:"openapi,omitempty"`
	// Paths    *map[string]interface{} `firestore:"openapi,omitempty"`
}

// API stores an API.
type API struct {
	Author   string
	Name     string
	Image    string
	Versions []APIVersion
}

func registerAPI(env *api.Env, author string, template APIVersion) error {
	im := api.ImageName{
		Author:  author,
		Name:    template.Manifest.Name,
		Version: template.Version,
	}

	ctx := context.Background()

	// Get API document.
	templateDoc := env.DB.Collection("apis").Doc(im.GetID())

	// Update API metadata.
	_, err := templateDoc.Set(ctx, map[string]interface{}{
		"author": im.Author,
		"name":   im.Name,
		"image":  env.GetDockerImagePath(im),
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	// Retrieve api version document
	versDoc := templateDoc.Collection("versions").Doc(im.Version)
	vers, versDocErr := versDoc.Get(ctx)
	if versDocErr != nil && status.Code(versDocErr) != codes.NotFound {
		return versDocErr
	}

	if versDocErr == nil {
		// Don't update if that API version exists and is public.
		tplVersion := APIVersion{}
		if err := vers.DataTo(&tplVersion); err != nil {
			return err
		}
		if tplVersion.IsPublic {
			err := fmt.Errorf("version %s already exists and is public", im.Version)
			return api.StatusError{Code: 400, Err: err}
		}
	}

	// Write version.
	_, errS := versDoc.Set(ctx, template)
	if errS != nil {
		return errS
	}

	return nil
}

// GetVersion retrieves an api version from the DB.
func GetVersion(env *api.Env, im api.ImageName) (*APIVersion, error) {
	// TODO: if "latest" retrieve all version of author/apiname and pick latest public.
	if im.Version == "latest" {
		return nil, errors.New("getting image tag 'latest' is not yet supported")
	}

	// Get API document.
	apiDoc := env.DB.Collection("apis").Doc(im.GetID())

	// Retrieve api version document.
	versDoc := apiDoc.Collection("versions").Doc(im.Version)
	vers, err := versDoc.Get(context.Background())
	if err != nil {
		return nil, err
	}

	// Retrieve Data.
	apiVersion := APIVersion{}
	if err := vers.DataTo(&apiVersion); err != nil {
		return nil, err
	}

	return &apiVersion, nil
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
