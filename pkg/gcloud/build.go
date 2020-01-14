package gcloud

import (
	"context"

	v1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
)

// CloudBuild implements ContainerBuilder for Cloud Build.
type CloudBuild struct {
	GCPProjectID string
}

// NewCloudBuild returns a new instance of Firestore.
func NewCloudBuild(ctx context.Context, gcpProjectID string) (*CloudBuild, error) {
	return &CloudBuild{
		GCPProjectID: gcpProjectID,
	}, nil
}

// GetBuildType returns the build type for skaffold.
func (cb *CloudBuild) GetBuildType() *v1.BuildType {
	return &v1.BuildType{
		GoogleCloudBuild: &v1.GoogleCloudBuild{
			ProjectID: cb.GCPProjectID,
		},
	}
}
