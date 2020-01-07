This is an [Istio Mixer Adapter](https://istio.io/docs/reference/config/policy-and-telemetry/mixer-overview/)
that adds basic API Management features to [Knative](https://knative.dev) using
[Firestore](https://firebase.com) as backend.

### Istio Mixer Adapter

One of the benefit of handling API management in an Istio Adapter is that
authorization is performed before the request hits Knative's activator & autoscaler.

So denied requests won't trigger pod autoscaling, and API management is entirely decoupled from
the microservice business logic and code.

![Istio Knative Mixer](https://istio.io/blog/2019/knative-activator-adapter/knative-mixer-adapter.png)
Diagram from [Istio's blog](https://istio.io/blog/2019/knative-activator-adapter/) of the Knative activator implemented as a mixer adapter.
In our case, we don't replace the Knative activator, we simply block unauthorized requests before they reach the activator.

![Istio Mixer Model Schema](https://raw.githubusercontent.com/wso2/istio-apim/master/request_flow.png)
Diagram from [WSO2 API Manager](https://github.com/wso2/istio-apim) illustrating how the adapter fits in the Istio Mixer Model.
In our case, the *API Manager Deployment* is just the Firestore.

### Firestore

To store & update quotas & metrics, this adapter manages 2 collections in firestore:

- `keys`: Each document represent 1 _consumer_ API key. Manually add a key to allow access to your API.
  Each key must have a `quotas` property with a `number` value `> 0`
  to successfuly access the endpoints.
- `requests`: contains 1 document per request with details useful for monitoring
  (timestamps, code..etc)

## Install

### Create service account & secret.

```bash
export KUDA_GCP_PROJECT="your-gcp-project-id"
export KUDA_ADAPTER_SA=kuda-adapter
export KUDA_ADAPTER_SA_EMAIL=$KUDA_ADAPTER_SA@$KUDA_GCP_PROJECT.iam.gserviceaccount.com

# Create the service account.
gcloud --project $KUDA_GCP_PROJECT iam service-accounts \
      create $KUDA_ADAPTER_SA \
      --display-name "Service Account for the deployer."

# Add firebase admin role
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_ADAPTER_SA_EMAIL \
  --role roles/firebase.admin

# Make a temporary directory to store key
KEY_DIRECTORY=$(mktemp -d)

# Download the secret key file for your service account.
gcloud iam service-accounts keys create $KEY_DIRECTORY/adapter-credentials.json \
  --iam-account=$KUDA_ADAPTER_SA_EMAIL

# Upload that as a secret in your Kubernetes cluster.
kubectl create secret -n istio-system generic adapter-credentials \
  --from-file=key.json=$KEY_DIRECTORY/adapter-credentials.json

# Delete the local secret
rm -rf $KEY_DIRECTORY
```

Deploy to kubernetes:

```bash
kubectl apply -f install/attributes.yaml
kubectl apply -f install/auth-template.yaml
kubectl apply -f install/metric-template.yaml
kubectl apply -f install/kuda.yaml
kubectl apply -f install/kuda_cfg.yaml
```

At this point, all requests to the `public` namespace should fail since we
haven't deployed the adapter (see Dev or Deploy sections).

## Deploy

```bash
kubectl apply -f install/adapter.yaml
```

## Dev

```bash
skaffold dev -f install/skaffold.yaml
```

Skaffold should automatically build, push & deploy the adapter's image to the
kubernetes cluster.

At this point, if you try to make a request to an endpoint in the `public`
namespace without a valid API key in the `x-api-key`, the adapter should reject
the request with an appropriate HTTP code (400, 403 or 429).

Example of a valid request:

```
curl -H 'x-api-key: XXXXXXXXX' "http://{YOUR_DOMAIN}"
```

If successful, the `credits_left` associated to this particular API key in Firestore should
have be decreased by 1. And there should be a new document with a few metrics
in the `requests` collection.

## Inspect / Debug:

#### Get Istio's mixer logs:

1. Get the mixer's pod id

```bash
kubectl -n istio-system get pods -lapp=policy
```

2. Stream logs from the mixer:

```bash
kubectl -n istio-system logs -f --tail 20 istio-policy-59dd46fb7d-948pt mixer
```

#### Mock requests locally:

Before starting the local mixer, you need to temporarily edit the file
`install/kuda_cfg.yaml`.

1. Comment the line 11 and uncomment the line 10: `address: "[::]:44225"`
2. Comment the line 60: `match: match(destination...etc`

Start local mixer server:

```bash
$ISTIO_BIN/mixs server \
    --configStoreURL=install \
    --log_output_level=attributes:debug
```

Start the adapter:

```bash
docker run --rm \
    -p 44225:44225 \
    -e FIRESTORE_CREDENTIALS="/path/YOUR_CREDENTIALS.json" \
    -v /path/to/local/credentials/folder:/secret \
    gcr.io/kuda-project/kuda-mixer-adapter
```

Send a mock Authorization request:

```bash
$ISTIO_BIN/mixc check \
    -s destination.service="svc.cluster.local" \
    -t request.time="2018-01-29T12:12:14Z",response.time="2018-01-29T12:12:15Z" \
    --stringmap_attributes "request.headers=x-api-key:ABCDxyz"
```

Send a mock Metric request:

```bash
$ISTIO_BIN/mixc report \
    -s destination.service="svc.cluster.local" \
    -t request.time="2018-01-29T12:12:14Z",response.time="2018-01-29T12:12:15Z" \
    --stringmap_attributes "request.headers=x-api-key:ABCDxyz"
```

## References:

- https://github.com/istio/istio/wiki/Mixer-Out-of-Process-Adapter-Walkthrough
- https://github.com/salrashid123/istio_custom_auth_adapter
- https://medium.com/@pubudu538/wso2-api-management-for-istio-service-mesh-6c682fc03835
- https://github.com/wso2/istio-apim
- https://github.com/apigee/istio-mixer-adapter
- https://istio.io/blog/2019/knative-activator-adapter/
