package gcloud

import (
	"context"
	"os"
	"os/exec"

	firebase "firebase.google.com/go"
	firebaseAuth "firebase.google.com/go/auth"
)

// FirebaseAuth implements api.Auth for Firebase authentication (Cloud IP)
type FirebaseAuth struct {
	Client *firebaseAuth.Client
}

// NewFirebaseAuth returns a new instance of FirebaseAuth.
func NewFirebaseAuth(ctx context.Context, gcpProjectID string) (*FirebaseAuth, error) {
	config := firebase.Config{ProjectID: gcpProjectID}
	app, err := firebase.NewApp(ctx, &config)
	if err != nil {
		return nil, err
	}
	auth, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &FirebaseAuth{
		Client: auth,
	}, nil
}

// VerifyIDToken verifies ID tokens and returns the token UID.
func (auth *FirebaseAuth) VerifyIDToken(accessToken string) (string, error) {
	token, err := auth.Client.VerifyIDToken(context.Background(), accessToken)
	if err != nil {
		return "", err
	}
	return token.UID, nil
}

// AuthServiceAccount authenticates gcloud using application credentials.
func AuthServiceAccount() error {
	cmd := exec.Command("gcloud", "auth", "activate-service-account", "--key-file",
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
