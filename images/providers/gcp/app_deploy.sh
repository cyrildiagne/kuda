#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

app_name=$(echo $1 | cut -f1 -d':')
app_version=$(echo $1 | cut -f2 -d':')
app_image="gcr.io/$KUDA_GCP_PROJECT_ID/$app_name:$app_version"
echo $app_image

# Build & push using cloud build.
gcloud builds submit --tag $app_image .

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Launch.
cat <<EOF | kubectl apply -f -
apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: $app_name
  namespace: default
spec:
  template:
    spec:
      containers:
        - image: $app_image
          resources:
            limits:
              nvidia.com/gpu: 1
EOF