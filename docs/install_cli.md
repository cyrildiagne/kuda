# Install CLI

### MacOS / Linux:

```bash
curl https://raw.githubusercontent.com/cyrildiagne/kuda/master/scripts/get-cli.sh -sSfL | sh
```

# Initialize

The CLI must be initialized with a remote Kuda cluster configuration.

<!--
## Using gpu.sh

The best way to get started quickly on a cost-effective, fully managed cluster.

First create an account on gpu.sh then initialize your local configuration with your namespace.

```bash
kuda init <your_namespace>
```

Replace <your_namespace> with your [gpu.sh](#) username.
-->


## Using [GCP](#)


GCP provides a good environment for running Kuda.

Follow the installation guide for
[installing Kuda on GCP](/docs/install_on_gcp.md).

Then initialize your local configuration with your namespace.

```bash
kuda init -n <your_namespace> -p <your_domain>
```