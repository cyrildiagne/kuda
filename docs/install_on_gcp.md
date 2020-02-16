# Install on GCP

This guide walks you through setting up Kuda on [GCP](https://cloud.google.com/kubernetes-engine/).

## Setup the GKE cluster

Requirements:

- [gcloud](#)
- [Kubectl](#)

First, make sure you've enabled the API services using gcloud:

```bash
gcloud services enable \
  cloudapis.googleapis.com \
  container.googleapis.com \
  containerregistry.googleapis.com \
  cloudbuild.googleapis.com
```

Then override some of the defaults settings to your configuration.
You can find the full list of config values in the [setup_gcp](/scripts/setup_gcp.sh) script.

```bash
export KUDA_GCP_PROJECT="your-gcp-project"
```

Finally run the `setup_gcp` script which will create the cluster
if it doesn't exist yet and will provision the required resources.

```bash
sh scripts/setup_gcp.sh
```

## Make sure Kubectl is connected to your cluster

```
kubectl get pods --all-namespaces
```

## Authentication

Install the authentication service, by following the instruction in
[/install/auth](/install/auth/README.md).

Then install Istio's authentication policy:

```
sh scripts/setup_auth_policy.sh
```

## API

Install the remote deployer API, by following the instructions in
[/install/api](/install/api).

## Enable HTTPS

You must have a real domain name (not xip.io auto-domain) to enable HTTPS.

The helper script enables HTTPS using [CloudDNS](#), [Let's Encrypt](#) and [cert-manager](#).
Adapt the ClusterIssuer manifest if you are using a different DNS.

```bash
export KUDA_DOMAIN="example.com"
export KUDA_NAMESPACE="default"
export KUDA_LETSENCRYPT_EMAIL="you@example.com"
sh scripts/gcp_enable_https.sh
```
