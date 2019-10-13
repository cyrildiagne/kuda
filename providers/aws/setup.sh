#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

function create_cluster() {
  clusterName=$1
  # cat <<EOF | eksctl create cluster --kubeconfig /aws-credentials/eksknative.yaml -f -

  cat <<EOF | eksctl create cluster -f -
    apiVersion: eksctl.io/v1alpha5
    kind: ClusterConfig

    metadata:
      name: $clusterName
      region: $KUDA_AWS_CLUSTER_REGION

    nodeGroups:
    - name: default
      instanceType: m5.large
      desiredCapacity: 2
    - name: gpu
      instanceType: p2.xlarge
      desiredCapacity: 1
EOF

}

function install_nvidia_drivers() {
  kubectl create \
    -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/1.0.0-beta3/nvidia-device-plugin.yml
}

function install_istio() {
  istio_folder="/istio-1.*"
  # Configure Helm
  if [ -z "$(kubectl -n kube-system get serviceaccounts | grep tiller)" ]; then
    kubectl create -f $istio_folder/install/kubernetes/helm/helm-service-account.yaml
    helm init --service-account tiller
  else
    echo "Helm/Tiller already installed."
  fi
  # Install prerequisites
  if [ -z "$(helm ls --all istio-init | grep istio-init)" ]; then
    helm install \
      --wait \
      --name istio-init \
      --namespace istio-system \
      $istio_folder/install/kubernetes/helm/istio-init
  else
    echo "Istio reprequisites already installed."
  fi
  # Install Istio
  echo "Installing Istio..."
  helm install \
    --wait \
    --name istio \
    --namespace istio-system \
    $istio_folder/install/kubernetes/helm/istio \
    --values $istio_folder/install/kubernetes/helm/istio/values-istio-demo.yaml
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

# Create cluster if it doesn't exist.
if [ -z "$(eksctl -v 0 get cluster --name $KUDA_AWS_CLUSTER_NAME --region $KUDA_AWS_CLUSTER_REGION)"]; then
  create_cluster $KUDA_AWS_CLUSTER_NAME
else
  echo "Cluster $KUDA_AWS_CLUSTER_NAME already exists."
fi

# Retrieve cluster token.
aws eks update-kubeconfig --name $KUDA_AWS_CLUSTER_NAME --region $KUDA_AWS_CLUSTER_REGION

# Install Nvidia drivers.
install_nvidia_drivers

# Install Istio.
if [ -z "$(helm ls --all istio | grep 'istio ')" ]; then
  install_istio
else
  echo "Istio is already installed."
fi

# Install Knative.
if [ -z "$(kubectl -n knative-serving | grep 'webhook')" ]; then
  install_knative
else
  echo "Knative is already installed."
fi

echo "Retrieving hostname:"
kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

echo "Cluster $KUDA_AWS_CLUSTER_NAME is ready!"
