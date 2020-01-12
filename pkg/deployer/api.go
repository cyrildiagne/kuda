package deployer

import (
	"context"
	"errors"
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
	im := ImageName{
		Author:  author,
		Name:    api.Manifest.Name,
		Version: api.Version,
	}

	ctx := context.Background()

	// Get API document.
	apiDoc := env.DB.Collection("apis").Doc(im.GetID())

	// Update API metadata.
	_, err := apiDoc.Set(ctx, map[string]interface{}{
		"author": im.Author,
		"name":   im.Name,
		"image":  env.GetDockerImagePath(im),
	}, firestore.MergeAll)
	if err != nil {
		return err
	}

	// Retrieve api version document
	versDoc := apiDoc.Collection("versions").Doc(im.Version)
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
			err := fmt.Errorf("version %s already exists and is public", im.Version)
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

// GetVersion retrieves an api version from the DB.
func GetVersion(env *Env, im ImageName) (*APIVersion, error) {
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
