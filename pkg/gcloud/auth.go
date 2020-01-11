package gcloud

import (
	"os"
	"os/exec"
)

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
