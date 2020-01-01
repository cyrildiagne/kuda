# Getting Started

Prerequisites:

- Kuda installed & cluster setup ([Installation Guide](install.md))
- A local copy of the [`hello-gpu-flask`](/examples/hello-gpu-flask) example:
  ```
  git clone github.com/cyrildiagne/kuda
  cd kuda/examples/hello-gpu-flask
  ```

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

To call & test your API, you need your cluster ingress IP address.
To retrieve it, you can run:

```bash
export cluster_ip=$(kubectl get svc istio-ingressgateway \
    --namespace istio-system \
    --output 'jsonpath={.status.loadBalancer.ingress[0].ip}')
```

You can then call the API, for example using cURL:

```bash
curl http://hello-gpu.default.$cluster_ip.xip.io
```

💡You can try to update the code in `main.py` while `dev`
is running and the remote API should automatically synchronize & reload
with the new changes.

## 4 - Deploy

Once you're happy with your API, you can deploy the production build using:

```bash
kuda deploy
```
