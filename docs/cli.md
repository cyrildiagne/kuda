# CLI Reference

## → Init

```bash
kuda init <deployer> [flags]
```

Initializes the local configuration.

**Arguments**

- `deployer` The API deployer. (default: `skaffold`)

**Flags**

- `[-n, --namespace]` Your namespace. (default: `default`)
- `[-d, --docker-registry]` Required when using the `skaffold` deployer.

**Examples**

```bash
kuda init -d gcr.io/my-gcp-project skaffold
```

<!--
```bash
kuda init \
    -n your-namespace
    -d gcr.io/my-gcp-project \
    localhost:8080
```

```bash
kuda init \
    -n your-namespace \
    deploy.kuda.gpu.sh
``` -->

## → Dev

```bash
kuda dev <manifest> [flags]
```

Deploys the API in development mode (with live file sync & app reload).

**Arguments**

- `manifest`·: Optional, The manifest file. (default: `kuda.yml`)

**Flags**

- `[--dry-run]` Generate the config files and skip execution.

**Examples**

```bash
kuda dev
```

```bash
kuda dev /path/to/manifest.yml
```

## → Deploy

```
kuda deploy <manifest> [flags]
```

Deploys the API in production mode.

**Arguments**

- `manifest`: Optional, The manifest file. (default: `kuda.yml`)

**Flags**

- `[--dry-run]` Generate the config files and skip deployment.

**Examples**

```bash
kuda deploy
```

```bash
kuda deploy /path/to/manifest.yml
```
