# Development

## Requirements

- go v1.13

## Setup

Clone the repository

## Implement a provider

A provider implementation should be a docker image tagged `gcr.io/kuda-project/provider-$KUDA_DEV_PROVIDER`. It should expose the following commands:

```bash
# Cluster
kuda_setup
kuda_delete
kuda_get

# Remote Dev Session
kuda_dev_start
kuda_dev_stop

# App
kuda_app_deploy
kuda_app_delete
```

To build a provider implementation & run a command:

Then:

```bash
make provider=<provider> cmd="<command>" build-provider-and-run
```

**Example:**

```bash
export KUDA_GCP_PROJECT_ID=gpu-sh
export KUDA_GCP_CREDENTIALS=~/Perso/kuda/secret/gpu-sh-f6d27675cda2.json
make provider=gcp cmd="get status" build-provider-and-run
```
