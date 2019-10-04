#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Delete dev pod.
kubectl delete service kuda-dev
kubectl delete deployment kuda-dev

# Delete istio resources
kubectl delete gateway kuda-dev-gateway
kubectl delete virtualservice kuda-dev

# Resize the GPU cluster to 0. > Not mandatory
# since the autoscaler will automatically scale down to 0 after a while.
# gcloud container clusters resize $KUDA_GCP_CLUSTER_NAME \
#   --node-pool $KUDA_DEFAULT_GPU \
#   --num-nodes 0 \
#   --quiet