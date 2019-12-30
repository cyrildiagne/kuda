package kuda

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
)

// CheckDeepEqual tests for deep equality
func CheckDeepEqual(t *testing.T, expected, actual interface{}) {
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

// GetTestConfig returns a test config
func GetTestConfig() Config {
	return Config{
		URLConfig: URLConfig{
			Scheme:    "https",
			Name:      "name",
			Namespace: "default",
			Domain:    "example.com",
		},
		DockerDestImage: "docker.io/test/test-api",
		Dockerfile:      "./Dockerfile",
		ManifestFile:    "./.kuda/service.yml",
		ConfigFolder:    "./.kuda",
	}
}

// GetTestDevConfig returns a test config
func GetTestDevConfig() Config {
	cfg := GetTestConfig()
	cfg.DevConfig = &DevConfig{
		Sync:    []string{"**/*.py"},
		Command: "cmd",
		Args:    []string{"a", "b", "c"},
		Env:     []corev1.EnvVar{{Name: "ENV_NAME", Value: "env-value"}},
	}
	return cfg
}
