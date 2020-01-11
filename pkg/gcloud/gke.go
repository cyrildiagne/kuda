package gcloud

import (
	"os"
	"os/exec"
)

// GetKubeConfig gets the kubeconfig.
func GetKubeConfig(gcpProjectID string) error {
	args := []string{"container", "clusters", "get-credentials",
		"--project", gcpProjectID,
		"--region", "us-central1-a", "kuda"}
	cmd := exec.Command("gcloud", args...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
