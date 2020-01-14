package gcloud

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/cyrildiagne/kuda/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Firestore implements api.DB for Firestore.
type Firestore struct {
	FS *firestore.Client
}

// NewFirestore returns a new instance of Firestore.
func NewFirestore(ctx context.Context, GCPProjectID string) (*Firestore, error) {
	client, err := firestore.NewClient(ctx, GCPProjectID)
	if err != nil {
		return nil, err
	}
	return &Firestore{
		FS: client,
	}, nil
}

// GetVersionnedAPI retrieves an api version from the DB.
func (fs *Firestore) GetVersionnedAPI(im api.ImageName) (*api.Version, error) {
	// TODO: if "latest" retrieve all version of author/apiname and pick latest public.
	if im.Version == "latest" {
		return nil, errors.New("getting image tag 'latest' is not yet supported")
	}

	// Get API document.
	apiDoc := fs.FS.Collection("apis").Doc(im.GetID())

	// Retrieve api version document.
	versDoc := apiDoc.Collection("versions").Doc(im.Version)
	vers, err := versDoc.Get(context.Background())
	if err != nil {
		return nil, err
	}

	// Retrieve Data.
	apiVersion := api.Version{}
	if err := vers.DataTo(&apiVersion); err != nil {
		return nil, err
	}

	return &apiVersion, nil
}

// IsUserAdminOfNamespace check if the namespace has UID as admin.
func (fs *Firestore) IsUserAdminOfNamespace(UID string, namespace string) (bool, error) {
	ctx := context.Background()
	// Retrieve the namespace
	ns, err := fs.FS.Collection("namespaces").Doc(namespace).Get(ctx)
	if err != nil {
		err = fmt.Errorf("error getting namespace info %v", err)
		return false, api.StatusError{Code: 500, Err: err}
	}
	if !ns.Exists() {
		err := fmt.Errorf("namespace not found %v", namespace)
		return false, api.StatusError{Code: 400, Err: err}
	}
	//
	nsData := ns.Data()
	nsAdmins, hasAdmins := nsData["admins"]
	if !hasAdmins {
		err := fmt.Errorf("no admin found for namespace %v", namespace)
		return false, api.StatusError{Code: 403, Err: err}
	}
	_, isAdmin := nsAdmins.(map[string]interface{})[UID]
	if !isAdmin {
		err := fmt.Errorf("user %v must be admin of %v", UID, namespace)
		return false, api.StatusError{Code: 403, Err: err}
	}
	return true, nil
}

// UpdateAPIMetadata updates API metadata
func (fs *Firestore) UpdateAPIMetadata(imageID string, metadata *map[string]interface{}) error {
	// Get API document.
	templateDoc := fs.FS.Collection("apis").Doc(imageID)

	// Update API metadata.
	ctx := context.Background()
	_, err := templateDoc.Set(ctx, *metadata, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// UpdateVersionnedAPI updates a versioned API.
func (fs *Firestore) UpdateVersionnedAPI(imID string, version string, template *api.Version) error {

	apiDoc := fs.FS.Collection("apis").Doc(imID)

	versDoc := apiDoc.Collection("versions").Doc(version)
	ctx := context.Background()
	vers, versDocErr := versDoc.Get(ctx)
	if versDocErr != nil && status.Code(versDocErr) != codes.NotFound {
		return versDocErr
	}

	if versDocErr == nil {
		// Don't update if that API version exists and is public.
		tplVersion := api.Version{}
		if err := vers.DataTo(&tplVersion); err != nil {
			return err
		}
		if tplVersion.IsPublic {
			err := fmt.Errorf("version %s already exists and is public", version)
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
