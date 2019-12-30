#!/bin/bash

# Requires gcloud and kubectl.
# Make sure you've enabled the API services using gcloud:
#   gcloud services enable \
#     cloudapis.googleapis.com \
#     container.googleapis.com \
#     containerregistry.googleapis.com

# Exit on error.
set -e

export PROJECT="${PROJECT:-gpu-sh}"
export DOMAIN="${DOMAIN:-xip.io}"
export NAMESPACE="${DOMAIN:-default}"
export CLUSTER_NAME="${CLUSTER_NAME:-kuda}"
export CLUSTER_ZONE="${CLUSTER_ZONE:-us-central1-a}"
export MASTER_MACHINE_TYPE="${MASTER_MACHINE_TYPE:-n1-standard-2}"
export GPU_MACHINE_TYPE="${GPU_MACHINE_TYPE:-n1-standard-2}"
export GPU_ACCELERATOR="${GPU_ACCELERATOR:-nvidia-tesla-k80}"
export USE_PREEMPTIBLE_GPU="${USE_PREEMPTIBLE_GPU:-true}"
export KNATIVE_VERSION="${KNATIVE_VERSION:-0.11.0}"

export CLUSTER_USER_ADMIN=$(gcloud config get-value core/account)

function create_main_cluster() {
  # Create the main Knative cluster.
  gcloud beta container clusters create $CLUSTER_NAME \
    --addons=HorizontalPodAutoscaling,HttpLoadBalancing,Istio \
    --machine-type=$MASTER_MACHINE_TYPE \
    --cluster-version=latest \
    --zone=$CLUSTER_ZONE \
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
  if [ $USE_PREEMPTIBLE_GPU = true ]; then
    preemptible_mode="--preemptible"
  fi
  # Create the default GPU Node pool.
  gcloud container node-pools create $GPU_ACCELERATOR \
    --machine-type=$GPU_MACHINE_TYPE \
    --accelerator type=$GPU_ACCELERATOR,count=1 \
    --zone $CLUSTER_ZONE \
    --cluster $CLUSTER_NAME \
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

  gcloud config set project $PROJECT

  # Check if cluster already exists otherwise create one.
  if gcloud container clusters list | grep -q $CLUSTER_NAME; then
    echo "→ Cluster already exists."
  else
    echo "Creating cluster $CLUSTER_NAME..."
    create_main_cluster
    grant_admin
    echo "→ Cluster created."
  fi

  # Get cluster's credentials to use kubectl.
  gcloud container clusters get-credentials $CLUSTER_NAME

  # Check if GPU cluster exists otherwise create one.
  if gcloud container node-pools list \
    --zone $CLUSTER_ZONE \
    --cluster $CLUSTER_NAME | grep -q $GPU_ACCELERATOR; then
    echo "→ GPU node pool already exists."
  else
    echo "Creating new GPU node pool with default GPU $GPU_ACCELERATOR..."
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
  if [ "$NAMESPACE" != "default" ]; then
    kubectl create namespace $NAMESPACE
  fi

  # Setup Domain name.
  if [ "$DOMAIN" = "xip.io" ]; then
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
      '{"data": {"example.com": null, "'$DOMAIN'": ""}}'
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
