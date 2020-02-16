## Deploy on GCP

### 1) Create service account and bind roles.

```bash
# Your GCP Project.
export KUDA_GCP_PROJECT="your-project-id"
# Name for the API service account that will be created.
export KUDA_API_SERVICE_ACCOUNT=kuda-api
# The full email for the service account.
export KUDA_API_SERVICE_ACCOUNT_EMAIL=$KUDA_API_SERVICE_ACCOUNT@$KUDA_GCP_PROJECT.iam.gserviceaccount.com

# Create the service account.
gcloud --project $KUDA_GCP_PROJECT iam service-accounts \
      create $KUDA_API_SERVICE_ACCOUNT \
      --display-name "Service Account for the deployer."

# Bind the role dns.admin to this service account, so it can be used to support
# the ACME DNS01 challenge.
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_API_SERVICE_ACCOUNT_EMAIL \
  --role roles/container.developer
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_API_SERVICE_ACCOUNT_EMAIL \
  --role roles/storage.objectCreator
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_API_SERVICE_ACCOUNT_EMAIL \
  --role roles/cloudbuild.builds.builder
gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
  --member serviceAccount:$KUDA_API_SERVICE_ACCOUNT_EMAIL \
  --role roles/firebase.admin
```

### 2) Create secret for this service account.

```bash
# Make a temporary directory to store key
KEY_DIRECTORY=$(mktemp -d)

# Download the secret key file for your service account.
gcloud iam service-accounts keys create $KEY_DIRECTORY/api-credentials.json \
  --iam-account=$KUDA_API_SERVICE_ACCOUNT_EMAIL

# Upload that as a secret in your Kubernetes cluster.
kubectl create secret -n kuda generic api-credentials \
  --from-file=key.json=$KEY_DIRECTORY/api-credentials.json

# Delete the local secret
rm -rf $KEY_DIRECTORY
```

### 3) Update the service.yaml with your GCP project id and project domain.

```bash
export KUDA_GCP_PROJECT="your-gcp-project"
export KUDA_DOMAIN="your-domain"
```

```bash
cd install/api
cp service-workaround.tpl.yaml service-workaround.yaml
sed -i'.bak' "s/\$KUDA_GCP_PROJECT/$KUDA_GCP_PROJECT/g" service-workaround.yaml
sed -i'.bak' "s/\$KUDA_DOMAIN/$KUDA_DOMAIN/g" service-workaround.yaml
rm service-workaround.yaml.bak
cd -
```

<!-- ```bash
cd install/api
cp service.tpl.yaml service.yaml
sed -i'.bak' "s/\$KUDA_GCP_PROJECT/$KUDA_GCP_PROJECT/g" service.yaml
sed -i'.bak' "s/\$KUDA_DOMAIN/$KUDA_DOMAIN/g" service.yaml
rm service.yaml.bak
cd - -->

### 4) Deploy.

```bash
kubectl apply -f install/api/service-workaround.yaml
```

<!-- ```bash
kubectl apply -f install/api/service.yaml
``` -->

Then check if your deployment is ready, `curl http://api.<your-domain>` and if
see "hello!", you are all set.

### 5) Create a Firestore Database.

If you're using the Firestore DB as database,
[create a database](https://console.cloud.google.com/firestore) for your project.

For now, you must also initialize the database with the namespaces you want to use.
So create a collection `namespaces`, add a document per namespace, each containing
a map `admins` of `<UserID> : <Timestamp` for each user allowed to deploy to this namespace.

For example:

```
COLLECTION   | DOCUMENTS   |
----------------------------------------------------------------------
namespaces > | default   > | > admins
             |             |      SEfsefsfBRsvsgXdqefs: February 2,...
```

The UserID of each user can be found in the identity provider (for instance GCloud Identity Platform).
Once they've logged in your service.

## Development

See [DEVELOPMENT.md](./DEVELOPMENT.MD)
