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
	cfg := NewConfig()
	expected := Config{
		ConfigFolder: "./.kuda",
		Dockerfile:   "./Dockerfile",
		KserviceFile: "./.kuda/service.yml",
		SkaffoldFile: "./.kuda/skaffold.yml",
	}
	if diff := cmp.Diff(expected, cfg); diff != "" {
		t.Errorf("TestNewConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestNewURLConfig(t *testing.T) {
	cfg := NewURLConfig()
	expected := URLConfig{
		Protocol:  "http",
		Namespace: "default",
	}
	if diff := cmp.Diff(expected, cfg); diff != "" {
		t.Errorf("TestNewURLConfig() mismatch (-want +got):\n%s", diff)
	}
}

func TestNewDevConfigFlask(t *testing.T) {
	cfg := NewDevConfigFlask()
	expected := DevConfig{
		Sync:    []string{"**/*.py"},
		Command: "python3",
		Args:    []string{"main.py"},
		Env: []corev1.EnvVar{
			{Name: "FLASK_ENV", Value: "development"},
		},
	}
	if diff := cmp.Diff(expected, cfg); diff != "" {
		t.Errorf("TestNewDevConfigFlask() mismatch (-want +got):\n%s", diff)
	}
}

func TestParseURL(t *testing.T) {
	emptyURL := ""

	if result, err := ParseURL(emptyURL); result != nil || err == nil {
		t.Errorf("ParseURL should return error with empty URL")
	}

	wrongURL := "../dir/"
	if result, err := ParseURL(wrongURL); result != nil || err == nil {
		t.Errorf("ParseURL should return error with wrongURL: %s", wrongURL)
	}

	urlA := "https://test-url.default.example.com"
	expectedA := URLConfig{
		Protocol:  "https",
		Name:      "test-url",
		Namespace: "default",
		Domain:    "example.com",
	}
	resultA, err := ParseURL(urlA)
	assert.NilError(t, err)
	if diff := cmp.Diff(expectedA, *resultA); diff != "" {
		t.Errorf("TestParseURL() mismatch (-want +got):\n%s", diff)
	}

	urlB := "http://test-url2.default.1.2.3.4.xip.io/run"
	expectedB := URLConfig{
		Protocol:  "http",
		Name:      "test-url2",
		Namespace: "default",
		Domain:    "1.2.3.4.xip.io",
	}
	resultB, err := ParseURL(urlB)
	assert.NilError(t, err)
	// if !reflect.DeepEqual(*resultB, expectedB) {
	// 	t.Errorf("Result B error. Got: \n%v, \nExpected: \n%v", resultB, expectedB)
	// }
	if diff := cmp.Diff(expectedB, *resultB); diff != "" {
		t.Errorf("TestParseURL() mismatch (-want +got):\n%s", diff)
	}
}

func TestGenerateConfigFiles(t *testing.T) {
	validURL := "https://test-url.default.example.com"
	validRegistry := "docker.io/test"
	manifests, err := GenerateConfigFiles(validURL, validRegistry)
	if err != nil {
		t.Errorf("Error generating valid config")
	}
	assert.Assert(t, len(manifests.Prod.Kservice) > 0)
	assert.Assert(t, len(manifests.Prod.Skaffold) > 0)
	assert.Assert(t, len(manifests.Dev.Kservice) > 0)
	assert.Assert(t, len(manifests.Dev.Skaffold) > 0)

	ksvc := knativev1.Service{}
	err = yaml.Unmarshal([]byte(manifests.Prod.Kservice), &ksvc)
	assert.NilError(t, err)
	err = yaml.Unmarshal([]byte(manifests.Dev.Kservice), &ksvc)
	assert.NilError(t, err)

	skfd := skaffoldv1.SkaffoldConfig{}
	err = yaml.Unmarshal([]byte(manifests.Prod.Skaffold), &skfd)
	assert.NilError(t, err)
	err = yaml.Unmarshal([]byte(manifests.Dev.Skaffold), &skfd)
	assert.NilError(t, err)
}
