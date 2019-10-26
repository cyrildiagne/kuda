#!/bin/bash

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
aws eks update-kubeconfig \
  --name $KUDA_AWS_CLUSTER_NAME \
  --region $KUDA_AWS_CLUSTER_REGION

# Apply the config file to the cluster.
kubectl apply -f $1