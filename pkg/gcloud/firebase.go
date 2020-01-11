package gcloud

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	firebaseAuth "firebase.google.com/go/auth"
)

// InitFirebase returns a firebase auth and firestore objects.
func InitFirebase(gcpProjectID string) (*firebaseAuth.Client, *firestore.Client, error) {
	config := &firebase.Config{ProjectID: gcpProjectID}
	app, err := firebase.NewApp(context.Background(), config)
	if err != nil {
		return nil, nil, err
	}

	auth, err := app.Auth(context.Background())
	if err != nil {
		return nil, nil, err
	}

	fs, err := app.Firestore(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return auth, fs, nil
}
