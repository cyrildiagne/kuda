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

## APP

### → Dev

```bash
kuda app dev <app-name>
```

Example: `kuda app dev my-app`

Starts the application in dev mode.

### → Deploy

```bash
kuda app deploy <app-name:version>
```

Example: `kuda app deploy my-app:1.0.0`

Deploys the application as serverless API.

### → Delete

```bash
kuda app delete <app-name>
```

Example: `kuda app delete my-app`

Deletes the application from the cluster and container image from the registry.
