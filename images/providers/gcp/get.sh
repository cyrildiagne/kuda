#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

function prop() {
  echo -e "-\e[1m $1 : \e[36m $2 \e[0m"
}

function get_status() {
  echo -e "\n\e[1mStatus:"
  # Print the IP Adress of the cluster.
  APP_IP_ADDRESS=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
  prop "IP" $APP_IP_ADDRESS
}

case $1 in
status) get_status ;;
*)
  echo "ERROR: command $1 unknown."
  exit
  ;;
esac
