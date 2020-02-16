# Getting Started

Prerequisites:

- Kuda installed & cluster setup ([Installation Guide](install.md))
- Kuda initialized to point to the cluster `kuda init <your_namespace> -p <your_domain>`


First, let's get a copy of the [`hello-gpu-flask`](/examples/hello-gpu-flask) example:
```
git clone github.com/cyrildiagne/kuda
cd kuda/examples/hello-gpu-flask
```

## 1 - Dev

Run the API remotely in dev mode using:

```bash
kuda dev
```

Depending on your configuration, the whole process could take a while.
But once the image has been built, pushed, provisioned, deployed & started,
you should start seeing the startup logs from the Flask debug server.

## 2 - Test

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

## 3 - Deploy

Once you're happy with your API, you can deploy the production build using:

```bash
kuda deploy
```

And call the production API, for example using cURL:

```bash
curl http://hello-gpu.default.your-domain.com
```