# CLI Reference

## → Init

```bash
kuda init <name> [-d, --docker-artifact] [-n, --namespace]
```

Generate the configuration files in a local `.kuda` folder.

**Arguments:**

- **`name`**: The API name. Must be unique in the namespace.

**Flags:**

- **`[-d | --docker-artifact]`**: A docker registry where you have write access (eg: `gcr.io/<your-gcp-project>/hello-gpu`).
- **`[-n | --namespace]`**: The Knative namespace (default: `default`).

## → Dev

```bash
kuda dev
```

Example: `kuda dev`

Deploys the API in dev mode.

## → Deploy 

```bash
kuda deploy
```

Example: `kuda deploy`

Deploys the API in production mode.
