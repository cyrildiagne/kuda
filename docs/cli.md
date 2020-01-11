# CLI Reference

## → Init

```bash
kuda init <namespace> [flags]
```

Initializes the local configuration.

**Arguments**

- `namespace` Your namespace.

**Flags**

- `[-p, --provider]` The provider root URL (default: `kuda.gpu.sh`).
- `[--auth-url]` Specify which url to use for authentication when using a remote deployer. (default: `auth.<provider>`)
- `[--deployer-url]` Specify which url to use for deployment when using a remote deployer. (default: `deployer.<provider>`)

**Examples**

```bash
kuda init my-namespace
```

```bash
kuda init my-namespace \
    -p kuda.my-cluster.com
```

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
