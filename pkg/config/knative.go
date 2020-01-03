package config

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "knative.dev/serving/pkg/apis/serving/v1"

	latest "github.com/cyrildiagne/kuda/pkg/manifest/latest"

	yaml "sigs.k8s.io/yaml"
)

// MarshalKnativeConfig generate yaml bytes from a knative config.
func MarshalKnativeConfig(s v1.Service) ([]byte, error) {
	content, err := yaml.Marshal(s)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GenerateKnativeConfig generate knative yaml specifics to the Kuda workflow.
func GenerateKnativeConfig(service ServiceSummary, cfg latest.Config) (v1.Service, error) {

	numGPUs, _ := resource.ParseQuantity("0")

	container := corev1.Container{
		Image: service.DockerArtifact,
		Name:  service.Name,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceName("nvidia.com/gpu"): numGPUs,
			},
		},
	}

	if cfg.Entrypoint.Command != "" {
		container.Command = []string{cfg.Entrypoint.Command}
	}
	if cfg.Entrypoint.Args != nil {
		container.Args = cfg.Entrypoint.Args
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
			Name:      service.Name,
			Namespace: service.Namespace,
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
