package gcloud

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// ListImageTags lists all tags of an image on gcr.io.
func ListImageTags(repoName string) error {
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
