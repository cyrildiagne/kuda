package deployer

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/cyrildiagne/kuda/pkg/manifest/latest"
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

func registerAPI(env *Env, author string, api APIVersion) error {
	name := api.Manifest.Name
	version := api.Version

	fullapiname := author + "__" + name

	ctx := context.Background()

	// Get API document.
	apiDoc := env.DB.Collection("apis").Doc(fullapiname)

	// Update API metadata.
	_, err := apiDoc.Set(ctx, map[string]interface{}{
		"author": author,
		"name":   name,
		"image":  env.GetDockerImagePath(author, name),
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	// Retrieve api version document
	versDoc := apiDoc.Collection("versions").Doc(version)
	vers, versDocErr := versDoc.Get(ctx)
	if versDocErr != nil && status.Code(versDocErr) != codes.NotFound {
		return versDocErr
	}

	if versDocErr == nil {
		// Don't update if that API version exists and is public.
		apiVersion := APIVersion{}
		if err := vers.DataTo(&apiVersion); err != nil {
			return err
		}
		if apiVersion.IsPublic {
			err := fmt.Errorf("version %s already exists and is public", version)
			return StatusError{400, err}
		}
	}

	// Write version.
	_, errS := versDoc.Set(ctx, api)
	if errS != nil {
		return errS
	}

	return nil
}
