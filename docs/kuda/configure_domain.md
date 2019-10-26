# Domain Configuration

## Setup the domain on the cluster

Create a configuration file and replace `your-domain.com` by your domain name:

```yaml
apiVersion: v1
   kind: ConfigMap
   metadata:
     name: config-domain
     namespace: knative-serving
   data:
     your-domain.com: ""
```

You can give this yaml configration file any name,
for example `domain-config.yaml`, then run:

```
kuda apply -f domain-config.yaml
```

You can see more options and informations about Knative's domain configuration
on [Knative's documentation](https://knative.dev/docs/serving/using-a-custom-domain/)

## Point the DNS domain to your cluster IP

You can then update your DNS provider to point the domain to the
IP adress of your cluster ingress.

To find out your cluster IP, you can run:
```bash
kubectl get svc istio-ingressgateway \
  --namespace istio-system \
  --output jsonpath="{.status.loadBalancer.ingress[*]['ip']}"
```

Find more information about setting custom domains
in [Knative's documentation](https://knative.dev/docs/serving/using-a-custom-domain/)