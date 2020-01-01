# Hello GPU HTTP API

A simple GPU HTTP API that just returns the output of `nvidia-smi` using Flask.

# Run locally

Requirements:

- [nvidia-docker](#)
- An Nvidia GPU compatible with CUDA

### 1 - Build the Docker image:

```bash
docker build hello-gpu .
```

### 2 - Start the server:

```bash
docker --runtime=nvidia-docker -p 8080:80 --rm run hello-gpu
```

### 3 - Test the API, for example using cURL:

```
curl http://localhost:8080
```

# Run as remote serverless API

Requirements:

- [Kuda CLI](#) pointing to a Kubernetes cluster with [Kuda](#).

## 1 - Initialize

Generate the configuration files using `kuda init`:

```bash
kuda init \
   -d docker.io/username/hello-gpu \
   hello-gpu
```

Replace `docker.io/username/hello-gpu` with a docker registry you have write
access to.

## 2 - Dev

Run the API remotely in dev mode using:

```bash
kuda dev
```

Depending on your configuration, the whole process could take a while.
But once the image has been built, pushed, provisioned, deployed & started,
you should start seeing the startup logs from the Flask debug server.

## 3 - Test

You can then call and test your dev API, for example using cURL (replace `your-domain.com` by your domain):

<details><summary>If you're using the (default) automatic xip.io domain</summary>
Then you need your cluster's ingress IP address to assemble the full domain name.
To retrieve it, you can run:

```bash
export cluster_ip=$(kubectl get svc istio-ingressgateway \
    --namespace istio-system \
    --output 'jsonpath={.status.loadBalancer.ingress[0].ip}')
echo "Your full xip.io domain is: $cluster_ip.xip.io"
```
</details>

```bash
curl http://hello-gpu-dev.default.your-domain.com
```

💡You can try to update the code in `main.py` while `dev`
is running and the remote API should automatically synchronize & reload
with the new changes.

## 4 - Deploy

Once you're happy with your API, you can deploy the production build using:

```bash
kuda deploy
```

And call the production API, for example using cURL:

```bash
curl http://hello-gpu.default.$cluster_ip.xip.io
```