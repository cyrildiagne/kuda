## Build

```bash
docker build \
  -t gcr.io/kuda-cloud/deployer \
  -f images/deployer/Dockerfile \
  .
```

## Run

```bash
docker run --rm \
  -e KUDA_GCP_PROJECT=`gcloud config get-value project` \
  -e GOOGLE_APPLICATION_CREDENTIALS=/credentials/`basename $GOOGLE_APPLICATION_CREDENTIALS` \
  -v `dirname $GOOGLE_APPLICATION_CREDENTIALS`:/credentials \
  -e PORT=80 \
  -p 8080:80 \
  gcr.io/kuda-cloud/deployer
```

## Deploy

### 1) Create service account and bind roles.

```bash
export KUDA_GCP_PROJECT="your-project-id"
export KUDA_DEPLOYER_SA=kuda-deployer
export KUDA_DEPLOYER_SA_EMAIL=$KUDA_DEPLOYER_SA@$KUDA_GCP_PROJECT.iam.gserviceaccount.com

# Create the service account.
gcloud --project $KUDA_GCP_PROJECT iam service-accounts \
      create $KUDA_DEPLOYER_SA \
      --display-name "Service Account for the deployer."

# Bind the role dns.admin to this service account, so it can be used to support
# the ACME DNS01 challenge.
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_DEPLOYER_SA_EMAIL \
  --role roles/container.developer
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_DEPLOYER_SA_EMAIL \
  --role roles/storage.objectCreator
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_DEPLOYER_SA_EMAIL \
  --role roles/cloudbuild.builds.builder
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_DEPLOYER_SA_EMAIL \
  --role roles/firebase.admin
```

### 2) Create secret for this service account.

```bash
# Make a temporary directory to store key
KEY_DIRECTORY=$(mktemp -d)

# Download the secret key file for your service account.
gcloud iam service-accounts keys create $KEY_DIRECTORY/deployer-credentials.json \
  --iam-account=$KUDA_DEPLOYER_SA_EMAIL

# Upload that as a secret in your Kubernetes cluster.
kubectl create secret -n kuda generic deployer-credentials \
  --from-file=key.json=$KEY_DIRECTORY/deployer-credentials.json

# Delete the local secret
rm -rf $KEY_DIRECTORY
```

### 3) Update the service.yaml with your GCP project id.

```bash
cp service.tpl.yaml service.yaml
sed -i'.bak' "s/value: <your-project-id>/value: $KUDA_GCP_PROJECT/g" service.yaml
rm service.yaml.bak
```

### 4) Deploy with skaffold.

```bash
skaffold run -f images/deployer/skaffold.yaml 
```

### 5) (Optional) If you want to start dev mode.

```bash
skaffold dev \
  -f images/deployer/skaffold.yaml 
```