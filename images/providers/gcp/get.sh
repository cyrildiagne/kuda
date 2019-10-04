#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

function prop() {
  echo -e "\e[1m$1 : \e[34m $2 \e[0m"
}

# Print the IP Adress of the cluster.
APP_IP_ADDRESS=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
prop "IP" $APP_IP_ADDRESS
