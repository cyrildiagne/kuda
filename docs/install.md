# Install

## 1 → Install CLI

### Requirements:

- [Docker](https://docs.docker.com/install)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Skaffold](https://skaffold.dev)

### MacOS / Linux:

Download the right binary from the [latest release](https://github.com/cyrildiagne/kuda/releases/latest).

Then move the binary to you \$PATH:

```bash
sudo mv kuda /usr/local/bin
```

---

## 2 → Setup Cluster

For now, Kuda is being developped actively on [GKE](https://cloud.google.com/kubernetes-engine/) but with cross-compatibility in mind.
Future releases will include setup scripts for other providers.

### Setup on GKE

Requirements:
- [gcloud](#)
- [Kubectl](#)

Make sure you've enabled the API services using gcloud:
```bash
gcloud services enable \
  cloudapis.googleapis.com \
  container.googleapis.com \
  containerregistry.googleapis.com
```

Then override some of the defaults settings to your configuration. You can find the full list of config values in the [setup_gcp](hack/setup_gcp.sh) scripts.
```bash
export PROJECT="your-gcp-project"
```

Finally run the `setup_gcp` script which will create the cluster if it doesn't exist yet and will provision the required resources.
```bash
sh hack/setup_gcp.sh
```