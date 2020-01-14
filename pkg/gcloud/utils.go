package gcloud

import (
	"context"
	"fmt"

	"github.com/cyrildiagne/kuda/pkg/api"
)

// NewEnv returns an api.Env for GCP.
func NewEnv(ctx context.Context, gcpProjectID string) (*api.Env, error) {
	auth, err := NewFirebaseAuth(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}

	cb, err := NewCloudBuild(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing cloud build: %v", err)
	}

	registry, err := NewGCR(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing registry: %v", err)
	}

	db, err := NewFirestore(ctx, gcpProjectID)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %v", err)
	}

	return &api.Env{
		Auth:              auth,
		ContainerRegistry: registry,
		ContainerBuilder:  cb,
		DB:                db,
	}, nil
}
