<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

🧪 **Status:** early and experimental WIP

---

## Rapidly develop and deploy serverless APIs that need GPUs

Kuda deploys APIs as serverless containers using [Knative](https://knative.dev)
which means you can use any language and any framework.
All you need is a Dockerfile.

## A simple interface for the full serverless API development cycle

- `kuda dev` Deploy the API on remote GPUs in dev mode (with file sync & live reload)
- `kuda deploy` Deploy the API in production mode.
  It will be automatically scaled down to zero when there is no traffic,
  and back up when there are new requests.

## Intuitive configuration with full control

Each API is configured with a simple declarative manifest such as:

```yaml
kudaManifestVersion: v1alpha1

# Name of the API.
name: hello-gpu

# 'deploy' is the configuration used when running `kuda deploy`.
# It has sensible defaults but you can override all the properties.
deploy:
  dockerfile: ./Dockerfile

# 'dev' is the configuration used when running `kuda dev`.
# It inherits all properties from 'deploy' which you can override individually.
dev:
  # Use python3 to start the Flask debug server rather than gunicorn.
  entrypoint:
    command: python3
    args: ["main.py"]
  # Live sync all python files.
  sync:
    - "**/*.py"
  # Set FLASK_ENV to "development" to enable Flask debugger & live reload.
  env:
    - name: FLASK_ENV
      value: development
```

## Run your APIs anywhere Kubernetes is running

<!-- - [gpu.sh](#) - The best way to get started quickly on a cost-effective, fully-managed GPU cluster. -->

- [GCP](#) - Installation guide for [running Kuda on GCP](/docs/install_on_gcp.md).

## Get Started

- [Install](docs/install_cli.md)
- [Getting Started](docs/getting_started.md)
- [CLI Reference](docs/cli.md)
