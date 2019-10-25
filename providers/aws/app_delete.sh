#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
app_name=$(echo $1 | cut -f1 -d':')

# Retrieve cluster token.
aws eks update-kubeconfig --name $KUDA_AWS_CLUSTER_NAME --region $KUDA_AWS_CLUSTER_REGION

# # Delete the knative service.
kubectl delete ksvc $app_name

# # Delete the ECS repository.
aws ecr delete-repository \
  --repository-name $app_name \
  --region $KUDA_AWS_CLUSTER_REGION \
  --force
# Delete cache repository.
aws ecr delete-repository \
  --repository-name "$app_name-cache" \
  --region $KUDA_AWS_CLUSTER_REGION \
  --force
