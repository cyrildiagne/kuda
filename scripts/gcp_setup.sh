#!/bin/bash

# Requires gcloud and kubectl.
# Make sure you've enabled the API services using gcloud:
#   gcloud services enable \
#     cloudapis.googleapis.com \
#     container.googleapis.com \
#     containerregistry.googleapis.com

# Exit on error.
set -e


red="\033[31m"
reset="\033[0m"

function print_help_and_exit() {
  echo "
This script requires the KUDA_GCP_PROJECT environment variables to be set.

Example usage:
  export KUDA_GCP_PROJECT=your-gcp-project
  sh scripts/gcp_setup.sh
"
  exit 1
}

function assert_set() {
  var_name=$1
  var_value=$2
  if [ -z "$var_value" ]; then
    printf "${red}ERROR:${reset} Missing required env variable $var_name\n"
    print_help_and_exit
  else
    echo "Using $var_name: $var_value"
  fi
}

assert_set KUDA_GCP_PROJECT $KUDA_GCP_PROJECT

export KUDA_DOMAIN="${KUDA_DOMAIN:-xip.io}"
export KUDA_NAMESPACE="${KUDA_NAMESPACE:-default}"
export KUDA_CLUSTER_NAME="${KUDA_CLUSTER_NAME:-kuda}"
export KUDA_CLUSTER_ZONE="${KUDA_CLUSTER_ZONE:-us-central1-a}"
export KUDA_MASTER_MACHINE_TYPE="${KUDA_MASTER_MACHINE_TYPE:-n1-standard-2}"
export KUDA_GPU_MACHINE_TYPE="${KUDA_GPU_MACHINE_TYPE:-n1-standard-2}"
export KUDA_GPU_ACCELERATOR="${KUDA_GPU_ACCELERATOR:-nvidia-tesla-k80}"
export KUDA_USE_PREEMPTIBLE_GPU="${KUDA_USE_PREEMPTIBLE_GPU:-true}"

exit 0

# The Knative version supported by this version of Kuda.
# Changing it might lead to unexpected behaviors.
KNATIVE_VERSION=0.11.0

# The user name to give RBAC admin role on the cluster.
CLUSTER_USER_ADMIN=$(gcloud config get-value core/account)

function create_main_cluster() {
  # Create the main Knative cluster.
  gcloud beta container clusters create $KUDA_CLUSTER_NAME \
    --addons=HorizontalPodAutoscaling,HttpLoadBalancing,Istio \
    --machine-type=$KUDA_MASTER_MACHINE_TYPE \
    --cluster-version=latest \
    --zone=$KUDA_CLUSTER_ZONE \
    --enable-stackdriver-kubernetes \
    --enable-ip-alias \
    --enable-autoscaling \
    --num-nodes=2 \
    --min-nodes=1 \
    --max-nodes=8 \
    --enable-autorepair \
    --enable-autoupgrade \
    --scopes cloud-platform
}

function grant_admin() {
  # Grant cluster-admin permissions to the current user.
  kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole=cluster-admin \
    --user=$CLUSTER_USER_ADMIN
}

function create_gpu_nodepools() {
  preemptible_mode=""
  if [ $KUDA_USE_PREEMPTIBLE_GPU = true ]; then
    preemptible_mode="--preemptible"
  fi
  # Create the default GPU Node pool.
  gcloud container node-pools create $KUDA_GPU_ACCELERATOR \
    --machine-type=$KUDA_GPU_MACHINE_TYPE \
    --accelerator type=$KUDA_GPU_ACCELERATOR,count=1 \
    --zone $KUDA_CLUSTER_ZONE \
    --cluster $KUDA_CLUSTER_NAME \
    --num-nodes 1 \
    --min-nodes 0 \
    --max-nodes 8 \
    --enable-autoupgrade \
    --enable-autoscaling \
    --metadata disable-legacy-endpoints=false \
    $preemptible_mode
}

function install_nvidia_drivers() {
  # Ensure sure the gcloud nvidia drivers are installed.
  NVIDIA_DRIVER_REPO="GoogleCloudPlatform/container-engine-accelerators"
  NVIDIA_DRIVER_PATH="master/nvidia-driver-installer/cos/daemonset-preloaded.yaml"
  kubectl apply -f "https://raw.githubusercontent.com/$NVIDIA_DRIVER_REPO/$NVIDIA_DRIVER_PATH"
}

function install_knative() {
  # Install Knative components.
  # We don't install monitoring nor eventing to save some resourceses
  knative_serving_repo="https://github.com/knative/serving/releases/download"
  kubectl apply --wait=true --selector knative.dev/crd-install=true \
    --filename $knative_serving_repo/v$KNATIVE_VERSION/serving.yaml

  # Hack to make sure the CRD have been installed correctly.
  sleep 3

  # Complete installation.
  kubectl apply \
    --filename $knative_serving_repo/v$KNATIVE_VERSION/serving.yaml
}

function setup() {
  echo "Setup Knative cluster on GKE..."

  gcloud config set project $KUDA_GCP_PROJECT

  # Check if cluster already exists otherwise create one.
  if gcloud container clusters list | grep -q $KUDA_CLUSTER_NAME; then
    echo "→ Cluster already exists."
  else
    echo "Creating cluster $KUDA_CLUSTER_NAME..."
    create_main_cluster
    grant_admin
    echo "→ Cluster created."
  fi

  # Get cluster's credentials to use kubectl.
  gcloud container clusters get-credentials $KUDA_CLUSTER_NAME

  # Check if GPU cluster exists otherwise create one.
  if gcloud container node-pools list \
    --zone $KUDA_CLUSTER_ZONE \
    --cluster $KUDA_CLUSTER_NAME | grep -q $KUDA_GPU_ACCELERATOR; then
    echo "→ GPU node pool already exists."
  else
    echo "Creating new GPU node pool with default GPU $KUDA_GPU_ACCELERATOR..."
    create_gpu_nodepools
    install_nvidia_drivers
    echo "→ GPU node pool created."
  fi

  # Install Knative.
  if kubectl get pods \
    --namespace knative-serving \
    --label-columns=serving.knative.dev/release | grep -q v$KNATIVE_VERSION; then
    echo "→ Knative v$KNATIVE_VERSION is already installed."
  else
    echo "Installing Knative v$KNATIVE_VERSION..."
    install_knative
    echo "→ Knative installed."
  fi

  # Setup namespace
  if [ "$KUDA_NAMESPACE" != "default" ]; then
    kubectl create namespace $KUDA_NAMESPACE
  fi

  # Setup Domain name.
  if [ "$KUDA_DOMAIN" = "xip.io" ]; then
    # TODO: remove this when after next Knative release.
    EXTERNAL_IP=$(kubectl get svc istio-ingressgateway \
      --namespace istio-system \
      --output jsonpath="{.status.loadBalancer.ingress[*]['ip']}")
    kubectl patch configmap config-domain \
      --namespace knative-serving \
      --patch \
      '{"data": {"example.com": null, "'$EXTERNAL_IP'.xip.io": ""}}'
  else
    kubectl patch configmap config-domain \
      --namespace knative-serving \
      --patch \
      '{"data": {"example.com": null, "'$KUDA_DOMAIN'": ""}}'
  fi
  echo "Done!"
}

function print_ip() {
  IP_ADDRESS=$(kubectl get svc istio-ingressgateway \
    --namespace istio-system \
    --output 'jsonpath={.status.loadBalancer.ingress[0].ip}')
  echo "Cluster Ingress IP Adress: $IP_ADDRESS"
}

setup
print_ip
