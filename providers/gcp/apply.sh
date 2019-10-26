#!/bin/bash

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Apply the config file to the cluster.
kubectl apply -f $1