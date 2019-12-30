# CLI Reference

## → Init

```bash
kuda init <URL> [-d, --docker-registry]
```

Generate the configuration files in a local `.kuda` folder.

**Arguments:**

- **`URL`**: The endpoint of the API.
  Kuda uses informations from this URL to generate the configuration files (service name, scheme, domain, namespace..Etc). Examples for an `hello-gpu` API:
  - `http://hello-gpu.default.12.34.56.78.xip.io` : The default configuration with http, the automatic domain xip.io and the `default` namespace.
  - `https://hello-gpu.public.yourdomain.com` : https with TLS termination, custom namespace `public` and a custom domain.

**Flags:**

- **`[-d | --docker-registry]`**: A docker registry where you have write access (eg: `docker.io/<your-username>`).

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
