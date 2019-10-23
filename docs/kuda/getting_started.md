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

Retrieve a simple demo application:

```bash
git clone https://github.com/cyrildiagne/kuda-apps
cd kuda-apps/hello-gpu
```

Then start the example in dev mode. It will reload automatically when you make changes from your local machine:

```bash
kuda app dev my-hello-gpu
```

Wait for the app to build and launch. This might take a while if a new node needs
to be allocated.

You can then query your application using any program able to make an HTTP request.
Here is an example using cURL:

```bash
curl -i -H "Host: my-hello-gpu.default.example.com" http://<YOUR-CLUSTER-IP>
```

Press `Ctrl+C` to stop running the application.

## 3 - Deploy

You can then deploy the app as a serverless API. This will create an endpoint that scales down the GPU nodes to 0 when not used.

From your local terminal run:

```bash
kuda app deploy hello-world:0.1.0
```

→ For more information on the `kuda app deploy` command, check the [reference](https://docs.kuda.dev/kuda/cli#deploy).

## 4 - Call your API

You can then test your application by making a simple HTTP request to your cluster.
First retrieve the IP address of your cluster by running: `kuda get status`

```bash
curl -i -H "Host: my-hello-gpu.default.example.com" http://<YOUR-CLUSTER-IP>
```

The first call might need to spawn an instance which could take while. Subsequent calls should be a lot faster.

## 5 - Cleanup

### • Delete the app

To delete the app (the image on the registry and the knative service):

```bash
kuda app delete hello-world
```

### • Delete the cluster

The GPU nodes are setup to autoscale down to 0 when they're not in use. However, the system node will still incur charges on GCP.

To completely delete the cluster run:

```bash
kuda delete
```

→ For more information on the `kuda cluster delete` command, check the [reference](https://docs.kuda.dev/kuda/cli#delete).
