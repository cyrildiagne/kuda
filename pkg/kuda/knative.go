package kuda

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "knative.dev/serving/pkg/apis/serving/v1"

	yaml "sigs.k8s.io/yaml"
)

// GenerateKnativeConfigYAML generate yaml string.
func GenerateKnativeConfigYAML(cfg Config) (string, error) {
	config, _ := GenerateKnativeConfig(cfg)
	content, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// GenerateKnativeConfig generate knative yaml specifics to the Kuda workflow
// and based on the kuda.Config given as parameter.
func GenerateKnativeConfig(cfg Config) (v1.Service, error) {

	if !cfg.IsValid() {
		return v1.Service{}, fmt.Errorf("invalid config")
	}

	numGPUs, _ := resource.ParseQuantity("1")

	container := corev1.Container{
		Image: cfg.DockerDestImage,
		Name:  cfg.Name,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceName("nvidia.com/gpu"): numGPUs,
			},
		},
	}

	if cfg.Command != "" {
		container.Command = []string{cfg.Command}
	}
	if cfg.Args != nil {
		container.Args = cfg.Args
	}
	if cfg.Env != nil {
		container.Env = append(container.Env, cfg.Env...)
	}

	config := v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},
		Spec: v1.ServiceSpec{
			ConfigurationSpec: v1.ConfigurationSpec{
				Template: v1.RevisionTemplateSpec{
					Spec: v1.RevisionSpec{
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{container},
						},
					},
				},
			},
		},
	}

	return config, nil
}
