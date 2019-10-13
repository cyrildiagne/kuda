#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
app_name=$(echo $1 | cut -f1 -d':')
app_version=$(echo $1 | cut -f2 -d':')
app_registry="$aws_account_id.dkr.ecr.$KUDA_AWS_CLUSTER_REGION.amazonaws.com/$app_name"
app_image="$app_registry:$app_version"

echo $app_image

# Create Container Registry if it doesn't exists.
if [ -z "$(aws ecr describe-repositories --region $KUDA_AWS_CLUSTER_REGION --repository-name $app_name)" ]; then
  aws ecr create-repository \
    --repository-name $app_name \
    --region $KUDA_AWS_CLUSTER_REGION
else
  echo "Container Registry $app_registry already exists"
fi

#TODO: Build image.


# Login Docker.
aws ecr get-login --region $KUDA_AWS_CLUSTER_REGION --no-include-email | bash

#TODO: Push image.


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