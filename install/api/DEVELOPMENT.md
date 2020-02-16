## Dev Locally

### Build

```bash
docker build \
  -t gcr.io/kuda-project/api \
  -f install/api/Dockerfile \
  .
```

### Run

```bash
docker run --rm \
  -e KUDA_GCP_PROJECT=`gcloud config get-value project` \
  -e GOOGLE_APPLICATION_CREDENTIALS=/credentials/`basename $GOOGLE_APPLICATION_CREDENTIALS` \
  -v `dirname $GOOGLE_APPLICATION_CREDENTIALS`:/credentials \
  -e PORT=80 \
  -p 8080:80 \
  gcr.io/kuda-project/api
```

### Deploy

```bash
docker push gcr.io/kuda-project/api
```

## Dev in the cluster using skaffold.

```bash
skaffold dev -f install/api/skaffold.yaml
```

## Deploy to the cluster with skaffold.

```bash
skaffold run -f install/api/skaffold.yaml
```
