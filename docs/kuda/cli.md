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

- **`[-p | --project]`**: A GCP project ID.
- **`[-c | --credentials]`**: Path to a GCP credentials file.

### → Delete

```bash
kuda delete
```

Deletes the remote cluster.

## DEV

### → Start

```bash
kuda dev start
```

Starts a remote dev session with the current working directory.

This command:

- Provisions a node with GPU on the cluster & install the nvidia driver
- Starts a development pod based on the Deep Learning VM
- Synchronise the directory provided as parameter with the remote node

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
