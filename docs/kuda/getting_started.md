# Getting Started

This guide gets you started with Kuda on GCP \(Google Cloud Platform\).

Appart from the environment variables, the process should be identical for other providers.

## 1 - Setup

- [Install Kuda](https://docs.kuda.dev/kuda/install) on your local machine
- Make sure you have an existing GCP Project with at least [1 quota for "GPUs \(all regions\)"](https://console.cloud.google.com/iam-admin/quotas).
- Download an application service [credentials as JSON](https://console.cloud.google.com/apis/credentials/serviceaccountkey)
- Then setup a new remote cluster:

```bash
kuda setup gcp --project <your-gcp-project-id> --credentials <path/to/your/credentials/json>
```

This process can take a while since it will create a remote cluster on GKE and install all the required addons.

→ For more information on the `kuda setup` command, check the [reference](https://docs.kuda.dev/kuda/cli#setup).

## 2 - Develop

### • Initialize

Retrieve an example application:

```bash
git clone https://github.com/cyrildiagne/kuda
cd kuda/examples/hello-world
```

Install the example dependencies (feel free to create a virtualenv or a [remote dev session](https://docs.kuda.dev/kuda/remote_development)).

```bash
pip install -r requirements.txt
```

### • Run and Test

Then start the example in dev mode. It will reload automatically when you make changes from your local machine:

```bash
export PORT=80 && python app.py
```

Open `http://localhost` in a web browser to visit the app. Try making changes to the code and reload the page.

Press `Ctrl+C` to stop running the application.

## • Deploy

You can then deploy the app as a serverless API. This will create an endpoint that scales down the GPU nodes to 0 when not used.

From your local terminal run:

```bash
kuda app deploy hello-world:0.1.0
```

→ For more information on the `kuda app deploy` command, check the [reference](https://docs.kuda.dev/kuda/cli#deploy).

## 3 - Call your API

You can then test your application by making a simple HTTP request to your cluster.
First retrieve the IP address of your cluster by running: `kuda get status`

```bash
curl -H "Host: hello-world.example.com" http://<cluster-ip-address>
```

The first call might need to spawn an instance which could take while. Subsequent calls should be a lot faster.

## 4 - Cleanup

### • Delete the cluster

The GPU nodes are setup to autoscale down to 0 when they're not in use. However, the system node will still incur charges on GCP.

To completely delete the cluster run:

```bash
kuda delete
```

→ For more information on the `kuda cluster delete` command, check the [reference](https://docs.kuda.dev/kuda/cli#delete).
