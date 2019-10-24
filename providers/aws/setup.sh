#!/bin/bash

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
      minSize: 0
      iam:
        withAddonPolicies:
          autoScaler: true
      tags:
        k8s.io/cluster-autoscaler/node-template/taint/dedicated: nvidia.com/gpu=true
        k8s.io/cluster-autoscaler/node-template/label/nvidia.com/gpu: "true"
        k8s.io/cluster-autoscaler/enabled: "true"
      labels:
        lifecycle: Ec2Spot
        nvidia.com/gpu: "true"
        k8s.amazonaws.com/accelerator: nvidia-tesla
      taints:
        nvidia.com/gpu: "true:NoSchedule"
EOF

}

function install_nvidia_drivers() {
  nvidia_driver_version=1.0.0-beta4
  nvidia_driver_host="https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin"
  kubectl create \
    -f $nvidia_driver_host/$nvidia_driver_version/nvidia-device-plugin.yml
}

function install_istio() {
  # Configure Helm
  istio_folder="/istio-1.*"
  if [ -z "$(kubectl -n kube-system get serviceaccounts | grep tiller)" ]; then
    kubectl create -f $istio_folder/install/kubernetes/helm/helm-service-account.yaml
    helm init --service-account tiller
  else
    echo "Helm/Tiller already installed."
  fi

  # Wait for Tiller to be ready.
  kubectl -n kube-system wait \
    --for=condition=Ready pod -l name=tiller --timeout=300s

  # Install prerequisites.
  if [ -z "$(helm ls --all istio-init | grep istio-init)" ]; then
    helm install \
      --wait \
      --name istio-init \
      --namespace istio-system \
      $istio_folder/install/kubernetes/helm/istio-init

      # Dirty hack to let the pods install.
      sleep 15
  else
    echo "Istio prerequisites already installed."
  fi

  # Install Istio.
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
  knative_version=0.9.0
  kubectl apply --wait=true --selector knative.dev/crd-install=true \
    --filename https://github.com/knative/serving/releases/download/v$knative_version/serving.yaml \
    --filename https://github.com/knative/eventing/releases/download/v$knative_version/release.yaml \
    --filename https://github.com/knative/serving/releases/download/v$knative_version/monitoring.yaml ||
      \
      kubectl apply --wait=true --selector knative.dev/crd-install=true \
      --filename https://github.com/knative/serving/releases/download/v$knative_version/serving.yaml \
      --filename https://github.com/knative/eventing/releases/download/v$knative_version/release.yaml \
      --filename https://github.com/knative/serving/releases/download/v$knative_version/monitoring.yaml

  kubectl apply \
    --filename https://github.com/knative/serving/releases/download/v$knative_version/serving.yaml \
    --filename https://github.com/knative/eventing/releases/download/v$knative_version/release.yaml \
    --filename https://github.com/knative/serving/releases/download/v$knative_version/monitoring.yaml
}

# Create cluster if it doesn't exist.
if [ -z "$(eksctl -v 0 get cluster --name $KUDA_AWS_CLUSTER_NAME --region $KUDA_AWS_CLUSTER_REGION)" ]; then
  create_cluster $KUDA_AWS_CLUSTER_NAME
else
  echo "Cluster $KUDA_AWS_CLUSTER_NAME already exists."
fi

# Retrieve cluster token.
aws eks update-kubeconfig --name $KUDA_AWS_CLUSTER_NAME --region $KUDA_AWS_CLUSTER_REGION

# Install Nvidia drivers.
if [ -z "$(kubectl get daemonsets -n kube-system nvidia-device-plugin-daemonset)" ]; then
  install_nvidia_drivers
else
  echo "Nvidia drivers already installed."
fi

# Install Istio.
install_istio
# if [ -z "$(helm ls --all istio | grep 'istio ')" ]; then
#   install_istio
# else
#   echo "Istio is already installed."
# fi

# Install Knative.
if [ -z "$(kubectl -n knative-serving get pods | grep 'webhook')" ]; then
  install_knative
else
  echo "Knative is already installed."
fi

# Create credentials for skaffold.
# https://github.com/GoogleContainerTools/skaffold/issues/1719
if [ -z "$(kubectl get secret | grep aws-secret)" ]; then
  kubectl create secret generic aws-secret --from-file /aws-credentials/credentials
else
  echo "aws-secret already configured."
fi

# Setup credential helpers in cluster for kaniko.
if [ -z "$(kubectl get secret | grep docker-kaniko-secret)" ]; then
  aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
  ecr_domain="$aws_account_id.dkr.ecr.$KUDA_AWS_CLUSTER_REGION.amazonaws.com"
  tmp_config_file="/tmp/config.json"
  echo "{ \"credHelpers\": { \"$ecr_domain\": \"ecr-login\" }}" > $tmp_config_file
  kubectl create secret generic docker-kaniko-secret --from-file $tmp_config_file
  rm $tmp_config_file
else
  echo "docker-kaniko-secret already configured."
fi

# Install Cluster Autoscaler
# echo "Installing cluster autoscaler"
# kubectl apply -f /kuda_cmd/config/cluster-autoscaler.yaml

echo
echo "Hostname:"
kubectl -n istio-system get service istio-ingressgateway \
  -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
echo

echo
echo "Cluster $KUDA_AWS_CLUSTER_NAME is ready!"
