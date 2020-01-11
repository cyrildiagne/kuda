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
- `[--auth-url]` Specify which url to use for authentication when using a remote deployer.
- `[--deployer-url]` Specify which url to use for deployment when using a remote deployer.

**Examples**

```bash
kuda init -d gcr.io/my-gcp-project skaffold
```

<!--
```bash
kuda init \
    -n your-namespace \
    gpu.sh
```

```bash
kuda init \
    -n your-namespace
    -d gcr.io/my-gcp-project \
    --auth_url localhost:8070 \
    --deployer_url localhost:8090
    localhost:8080
```

-->

## → Dev

```bash
kuda dev
```

Deploys the API in development mode (with live file sync & app reload).

**Examples**

Deploy the API from the local directory:

```bash
kuda dev
```

## → Deploy

```
kuda deploy
```

Deploys the API in production mode.

**Flags**

- `[-f, --from]` Qualitifed name of a published API from the registry: `<user>/<api-name>:<version|latest>`

**Examples**

Deploy the API from the local directory:

```bash
kuda deploy
```

Deploy the API from a published API in the repo:

```bash
kuda deploy -f cyrildiagne/hello-gpu
kuda deploy -f cyrildiagne/hello-gpu:1.3.0
```

## → Publish

```
kuda publish
```

Publish the API template to the registry.
This command publishes the API template & docker image so that other users can
deploy it inside their own environment.
It doesn't affect access to your deployed APIs.

**Examples**

Publish the API template from the local directory:

```bash
kuda publish
```