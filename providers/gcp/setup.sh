#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# required_addons="HorizontalPodAutoscaling,HttpLoadBalancing,Istio,CloudRun"
required_addons="HorizontalPodAutoscaling,HttpLoadBalancing,Istio"

function create_main_cluster() {
  # Create cluster with the system node pool.
  gcloud beta container clusters create $KUDA_GCP_CLUSTER_NAME \
    --addons=$required_addons \
    --machine-type=n1-standard-4 \
    --cluster-version=latest \
    --zone=$KUDA_GCP_COMPUTE_ZONE \
    --scopes cloud-platform \
    --num-nodes $KUDA_GCP_POOL_NUM_NODES \
    --enable-stackdriver-kubernetes \
    --issue-client-certificate \
    --enable-basic-auth \
    --enable-ip-alias \
    --enable-autoupgrade \
    --metadata disable-legacy-endpoints=false

  # Grant cluster-admin permissions to the current user.
  kubectl create clusterrolebinding cluster-admin-binding \
    --clusterrole=cluster-admin \
    --user=$(gcloud config get-value core/account)
}

function create_gpu_nodepools() {
  preemptible_mode=""
  if [ $KUDA_GCP_USE_PREEMPTIBLE = true ]; then
    preemptible_mode='--preemptible'
  fi
  # Create the default GPU Node pool.
  gcloud container node-pools create $KUDA_GCP_GPU \
    --machine-type=$KUDA_GCP_MACHINE_TYPE \
    --accelerator type=nvidia-tesla-$KUDA_GCP_GPU,count=1 \
    --zone $KUDA_GCP_COMPUTE_ZONE \
    --cluster $KUDA_GCP_CLUSTER_NAME \
    --num-nodes 1 \
    --min-nodes 0 \
    --max-nodes 8 \
    --enable-autoupgrade \
    --enable-autoscaling \
    --metadata disable-legacy-endpoints=false \
    $preemptible_mode
}

function install_knative() {
  # Apply crd twice as workaround to https://github.com/knative/serving/issues/5722
  kubectl apply --wait=true --selector knative.dev/crd-install=true \
    --filename https://github.com/knative/serving/releases/download/v0.9.0/serving.yaml \
    --filename https://github.com/knative/eventing/releases/download/v0.9.0/release.yaml \
    --filename https://github.com/knative/serving/releases/download/v0.9.0/monitoring.yaml ||
      \
      kubectl apply --wait=true --selector knative.dev/crd-install=true \
      --filename https://github.com/knative/serving/releases/download/v0.9.0/serving.yaml \
      --filename https://github.com/knative/eventing/releases/download/v0.9.0/release.yaml \
      --filename https://github.com/knative/serving/releases/download/v0.9.0/monitoring.yaml

  kubectl apply \
    --filename https://github.com/knative/serving/releases/download/v0.9.0/serving.yaml \
    --filename https://github.com/knative/eventing/releases/download/v0.9.0/release.yaml \
    --filename https://github.com/knative/serving/releases/download/v0.9.0/monitoring.yaml
}

function enable_gcloud_services() {
  # Enable the required Gcloud services.
  echo "Checking Gcloud services availability..."
  gcloud_services=$(gcloud services list)
  if [ -z "$(echo $gcloud_services | grep serviceusage)" ]; then
    echo "Enabling serviceusage API"
    gcloud services enable serviceusage.googleapis.com
  fi
  if [ -z "$(echo $gcloud_services | grep cloudapis)" ]; then
    echo "Enabling serviceusage API"
    gcloud services enable cloudapis.googleapis.com
  fi
  if [ -z "$(echo $gcloud_services | grep container)" ]; then
    echo "Enabling container API"
    gcloud services enable container.googleapis.com
  fi
  if [ -z "$(echo $gcloud_services | grep containerregistry)" ]; then
    echo "Enabling container registry API"
    gcloud services enable containerregistry.googleapis.com
  fi
  if [ -z "$(echo $gcloud_services | grep cloudbuild)" ]; then
    echo "Enabling cloud build API"
    gcloud services enable cloudbuild.googleapis.com
  fi
}

enable_gcloud_services

# Check if cluster already exists otherwise create one.
if gcloud container clusters list | grep -q $KUDA_GCP_CLUSTER_NAME; then
  echo "Cluster already exists."
else
  echo "Creating new cluster $KUDA_GCP_CLUSTER_NAME"
  create_main_cluster
fi

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Check if GPU cluster exists otherwise create one.
if gcloud container node-pools list \
  --zone $KUDA_GCP_COMPUTE_ZONE \
  --cluster $KUDA_GCP_CLUSTER_NAME | grep -q $KUDA_GCP_GPU; then
  echo "GPU node pool already exists."
else
  echo "Creating new GPU node pool with default GPU $KUDA_GCP_GPU"
  create_gpu_nodepools
fi

# Make sure the Nvidia drivers are installed
nvidia_driver_repo="https://raw.githubusercontent.com/GoogleCloudPlatform/container-engine-accelerators"
nvidia_driver_path="master/nvidia-driver-installer/cos/daemonset-preloaded.yaml"
kubectl apply -f "$nvidia_driver_repo/$nvidia_driver_path"

# Install Knative.
if kubectl get pods --namespace knative-serving --label-columns=serving.knative.dev/release | grep v0.9.0; then
  echo "Knative v0.9.0 is already installed."
else
  echo "Installing Knative 0.9.0..."
  install_knative
fi

# Create namespaces.
# TODO: hide the output of this command (especially if no namespace found.)
if [ -z "$(kubectl get namespace kuda-app)" ]; then
  kubectl create namespace kuda-app
fi

# if [ -z "$(kubectl get namespace kuda-dev)" ]; then
#   kubectl create namespace kuda-dev
#   kubectl label namespace kuda-dev istio-injection=enabled
# fi

# Enable Istio sidecar injection.
# kubectl label namespace kuda-app istio-injection=enabled --overwrite=true

# Mount credentials.
if kubectl get secrets | grep -q $(basename $KUDA_GCP_CREDENTIALS); then
  echo "Secret already exists."
else
  kubectl create secret generic $(basename $KUDA_GCP_CREDENTIALS) \
    --from-file=$KUDA_GCP_CREDENTIALS
fi