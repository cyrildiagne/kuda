module github.com/cyrildiagne/kuda

go 1.13

replace (
	contrib.go.opencensus.io/exporter/stackdriver => contrib.go.opencensus.io/exporter/stackdriver v0.12.9-0.20191108183826-59d068f8d8ff
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.2
	github.com/docker/docker => github.com/docker/docker v1.4.2-0.20191212201129-5f9f41018e9d
	golang.org/x/crypto v0.0.0-20190129210102-0709b304e793 => golang.org/x/crypto v0.0.0-20180904163835-0709b304e793
)

require (
	github.com/GoogleContainerTools/skaffold v1.1.0
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/google/go-cmp v0.3.1
	github.com/gorilla/mux v1.7.3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	gopkg.in/yaml.v2 v2.2.7
	gotest.tools v2.2.0+incompatible
	k8s.io/api v0.0.0-20190831074750-7364b6bdad65
	k8s.io/apimachinery v0.0.0-20190831074630-461753078381
	knative.dev/pkg v0.0.0-20191230041935-400dfb9ff95a // indirect
	knative.dev/serving v0.11.1
	sigs.k8s.io/yaml v1.1.0
)
