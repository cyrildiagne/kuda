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
		Name:            "test",
		Namespace:       "test",
		DockerDestImage: "test.io/test/test",
		Dockerfile:      "./Dockerfile",
		KserviceFile:    "./.kuda/service.yml",
		SkaffoldFile:    "./.kuda/skaffold.yml",
		ConfigFolder:    "./.kuda",
	}
}

// GetTestDevConfig returns a test config
func GetTestDevConfig() Config {
	cfg := GetTestConfig()
	cfg.Sync = []string{"**/*.py"}
	cfg.Command = "cmd"
	cfg.Args = []string{"a", "b", "c"}
	cfg.Env = []corev1.EnvVar{{Name: "ENV_NAME", Value: "env-value"}}
	return cfg
}
