# Install CLI

Requirements:

- [Docker](https://docs.docker.com/install)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
- [Skaffold](https://skaffold.dev)

### MacOS / Linux:

```bash
curl https://raw.githubusercontent.com/cyrildiagne/kuda/master/scripts/get-kuda-cli.sh -sSfL | sh
```

---

# Setup Cluster

Kuda is being actively developped on [GKE](https://cloud.google.com/kubernetes-engine/)
but with cross-compatibility in mind.
Future releases will include setup scripts for other providers.

### Setup on GKE

Requirements:

- [gcloud](#)
- [Kubectl](#)

First, make sure you've enabled the API services using gcloud:

```bash
gcloud services enable \
  cloudapis.googleapis.com \
  container.googleapis.com \
  containerregistry.googleapis.com
```

Then override some of the defaults settings to your configuration.
You can find the full list of config values in the [setup_gcp](/scripts/setup_gcp.sh) script.

```bash
export PROJECT="your-gcp-project"
```

Finally run the `setup_gcp` script which will create the cluster
if it doesn't exist yet and will provision the required resources.

```bash
./scripts/setup_gcp.sh
```
