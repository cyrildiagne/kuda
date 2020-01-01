package kuda

import (
	"testing"

	skaffoldv1 "github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/v1"
	"github.com/google/go-cmp/cmp"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	knativev1 "knative.dev/serving/pkg/apis/serving/v1"
	yaml "sigs.k8s.io/yaml"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("test", "test")
	expected := Config{
		Name:         "test",
		Namespace:    "test",
		ConfigFolder: "./.kuda",
		Dockerfile:   "./Dockerfile",
		KserviceFile: "./.kuda/service.yml",
		SkaffoldFile: "./.kuda/skaffold.yml",
	}
	if diff := cmp.Diff(expected, cfg); diff != "" {
		t.Errorf("TestNewConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestAddDevConfigFlask(t *testing.T) {
	cfg := NewConfig("test", "test")
	cfg.AddDevConfigFlask()

	expected := NewConfig("test", "test")
	expected.Sync = []string{"**/*.py"}
	expected.Command = "python3"
	expected.Args = []string{"main.py"}
	expected.Env = []corev1.EnvVar{
		{Name: "FLASK_ENV", Value: "development"},
	}

	if diff := cmp.Diff(expected, cfg); diff != "" {
		t.Errorf("TestAddDevConfigFlask() mismatch (-want +got):\n%s", diff)
	}
}

func TestSetFilesSuffix(t *testing.T) {
	cfg := NewConfig("test", "test")
	cfg.SetFilesSuffix("-dev")
	assert.Equal(t, cfg.KserviceFile, "./.kuda/service-dev.yml")
	assert.Equal(t, cfg.SkaffoldFile, "./.kuda/skaffold-dev.yml")
}

func TestIsValid(t *testing.T) {
	cfg := Config{}
	assert.Assert(t, !cfg.IsValid())

	cfg = NewConfig("test", "test")
	assert.Assert(t, cfg.IsValid())

	cfg.ConfigFolder = ""
	assert.Assert(t, !cfg.IsValid())

	cfg = NewConfig("test", "test")
	cfg.Dockerfile = ""
	assert.Assert(t, !cfg.IsValid())

	cfg = NewConfig("test", "test")
	cfg.KserviceFile = ""
	assert.Assert(t, !cfg.IsValid())
}

func TestGenerateConfigFiles(t *testing.T) {
	cfgDev := NewConfig("test", "test")
	cfgDev.DockerDestImage = "test.io/test"
	manifests, err := GenerateConfigFiles(cfgDev)
	if err != nil {
		t.Errorf("Error generating valid config")
	}
	assert.Assert(t, len(manifests.Kservice) > 0)
	assert.Assert(t, len(manifests.Skaffold) > 0)

	ksvc := knativev1.Service{}
	err = yaml.Unmarshal([]byte(manifests.Kservice), &ksvc)
	assert.NilError(t, err)

	skfd := skaffoldv1.SkaffoldConfig{}
	err = yaml.Unmarshal([]byte(manifests.Skaffold), &skfd)
	assert.NilError(t, err)
}
