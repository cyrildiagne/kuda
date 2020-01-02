# Install CLI

### MacOS / Linux:

```bash
curl https://raw.githubusercontent.com/cyrildiagne/kuda/master/scripts/get-cli.sh -sSfL | sh
```

# Initialize

The CLI must be initialized with a remote Kuda cluster configuration.

<!-- ## Using gpu.sh

The best way to get started quickly on a cost-effective, fully managed cluster.

```bash
kuda init \
  -n $your_namespace \
  deploy.kuda.gpu.sh 
```

Replace `$your_namespace` with your [gpu.sh](#) username. -->

## Using [GCP](#)


GCP provides a good environment for running Kuda.

Follow the installation guide for
[installing Kuda on GCP](/docs/install_on_gcp.md).

Then configure the CLI to deploy to the GKE cluster directly from your
workstation using the `skaffold` deployer.
This deployer requires [Docker](docker.com) and [Skaffold](https://skaffold.dev)
installed and configured on your machine.

```bash
kuda init \
  -n $your_namespace \
  -d gcr.io/$your_gcp_project \
  skaffold
```

<!-- ```bash
If you've installed and configured a [Kuda Deployer](#): 

kuda init \
  -n $your_namespace \
  deploy.kuda.$your_domain
```
-->

