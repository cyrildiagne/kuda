#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

app_name=$(echo $1 | cut -f1 -d':')
app_version=$(echo $1 | cut -f2 -d':')
app_image="gcr.io/$KUDA_GCP_PROJECT_ID/$app_name:$app_version"
echo $app_image

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Delete Knative service.
kubectl delete -n kuda-app ksvc $app_name

# Delete image from the repository.
gcloud container images delete $app_image