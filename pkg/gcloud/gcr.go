package gcloud

import (
	"context"
	"fmt"

	"github.com/cyrildiagne/kuda/pkg/api"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// GCR implements ContainerRegistry for GCR.
type GCR struct {
	GCPProjectID string
}

// NewGCR returns a new instance of GCR.
func NewGCR(ctx context.Context, gcpProjectID string) (*GCR, error) {
	return &GCR{
		GCPProjectID: gcpProjectID,
	}, nil
}

// GetDockerImagePath returns the fully qualified URL of a docker image on
// the appropriate registry.
func (gcr *GCR) GetDockerImagePath(im api.ImageName) string {
	return "gcr.io/" + gcr.GCPProjectID + "/" + im.GetID()
}

// ListImageTags lists all tags of an image on gcr.io.
func (gcr *GCR) ListImageTags(repoName string) error {
	repo, err := name.NewRepository(repoName)
	if err != nil {
		return err
	}
	fmt.Println(repo.Name())
	tags, err := google.List(repo, google.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}
	fmt.Println(tags)
	return nil
}
