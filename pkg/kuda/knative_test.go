package kuda

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "knative.dev/serving/pkg/apis/serving/v1"
)

func TestGenerateKnativeConfig(t *testing.T) {

	cfg := GetTestConfig()

	result, err := GenerateKnativeConfig(cfg)
	if err != nil {
		t.Errorf("err")
	}

	CheckDeepEqual(t, result.APIVersion, "serving.knative.dev/v1")
	CheckDeepEqual(t, result.Kind, "Service")

	meta := metav1.ObjectMeta{
		Name:      cfg.URLConfig.Name,
		Namespace: cfg.URLConfig.Namespace,
	}
	CheckDeepEqual(t, result.ObjectMeta, meta)

	numGPUs, _ := resource.ParseQuantity("1")
	container := corev1.Container{
		Image: cfg.DockerDestImage,
		Name:  cfg.URLConfig.Name,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceName("nvidia.com/gpu"): numGPUs,
			},
		},
	}
	spec := v1.ServiceSpec{
		ConfigurationSpec: v1.ConfigurationSpec{
			Template: v1.RevisionTemplateSpec{
				Spec: v1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{container},
					},
				},
			},
		},
	}
	CheckDeepEqual(t, result.Spec, spec)
}

func TestGenerateKnativeDevConfig(t *testing.T) {

	cfg := GetTestDevConfig()

	result, err := GenerateKnativeConfig(cfg)
	if err != nil {
		t.Errorf("err")
	}

	numGPUs, _ := resource.ParseQuantity("1")
	container := corev1.Container{
		Image: cfg.DockerDestImage,
		Name:  cfg.URLConfig.Name,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceName("nvidia.com/gpu"): numGPUs,
			},
		},
		Command: []string{"cmd"},
		Args:    []string{"a", "b", "c"},
		Env:     []corev1.EnvVar{{Name: "ENV_NAME", Value: "env-value"}},
	}
	spec := v1.ServiceSpec{
		ConfigurationSpec: v1.ConfigurationSpec{
			Template: v1.RevisionTemplateSpec{
				Spec: v1.RevisionSpec{
					PodSpec: corev1.PodSpec{
						Containers: []corev1.Container{container},
					},
				},
			},
		},
	}
	CheckDeepEqual(t, result.Spec, spec)
}
