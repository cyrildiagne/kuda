# Reference

## Setup

### → Setup

```bash
kuda setup [gcp]
```

Setup a remote cluster on given provider.

- **`provider`**: A managed kubernetes provider. For now only GCP is implemented.

This command also allocates an empty GPU node that will be provisioned only when you launch remote sessions.

**Flags for `gcp`:**

- **`[-p | --gcp_project_id]`**: An existing GCP project ID.
- **`[-c | --gcp_credentials]`**: Path to a GCP credentials file.

### → Delete

```bash
kuda delete
```

Deletes the remote cluster.

## DEV

### → Start

```bash
kuda dev start [base-image]
```

Starts a remote dev session with the current working directory.

**Example:** `kuda dev start nvidia/cuda:10.1-base`

This command:

- Provisions a node with GPU on the cluster & install the nvidia driver
- Starts a development pod based on the Deep Learning VM
- Synchronise the directory provided as parameter with the remote node

List of recommended `base-image`:

- all images from [nvidia/cuda](https://hub.docker.com/r/nvidia/cuda/)
- gcloud's [Deep Learning containers](https://cloud.google.com/ai-platform/deep-learning-containers/docs/choosing-container)

### → Stop

```bash
kuda dev stop
```

Stops the remote dev session.

## APP

### → Deploy

```bash
kuda app deploy <app-name:version>
```

Example: `kuda app deploy my-app:1.0.0`

Deploys the application as serverless API.
