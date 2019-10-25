#!/bin/bash

source $KUDA_CMD_DIR/.config.sh

cluster_region=$KUDA_AWS_CLUSTER_REGION
cluster_name=$KUDA_AWS_CLUSTER_NAME

function create_cluster() {
  cat <<EOF | eksctl create cluster -f -
    apiVersion: eksctl.io/v1alpha5
    kind: ClusterConfig

    metadata:
      name: $cluster_name
      region: $cluster_region

    nodeGroups:
    - name: default
      instanceType: m5.large
      desiredCapacity: 2
    # - name: gpu
    #   instanceType: p2.xlarge
    #   desiredCapacity: 1
    #   minSize: 0
    #   # iam:
    #   #   withAddonPolicies:
    #   #     autoScaler: true
    #   # tags:
    #   #   k8s.io/cluster-autoscaler/node-template/taint/dedicated: nvidia.com/gpu=true
    #   #   k8s.io/cluster-autoscaler/node-template/label/nvidia.com/gpu: "true"
    #   #   k8s.io/cluster-autoscaler/enabled: "true"
    #   # labels:
    #   #   lifecycle: Ec2Spot
    #   #   nvidia.com/gpu: "true"
    #   #   k8s.amazonaws.com/accelerator: nvidia-tesla
    #   # taints:
    #   #   nvidia.com/gpu: "true:NoSchedule"
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
  is_tiller_installed="$(kubectl -n kube-system get serviceaccounts | grep tiller)"
  if [ -z "$is_tiller_installed" ]; then
    kubectl create \
      -f $istio_folder/install/kubernetes/helm/helm-service-account.yaml
    helm init --service-account tiller
  else
    echo "Helm/Tiller already installed."
  fi

  # Wait for Tiller to be ready.
  kubectl -n kube-system wait \
    --for=condition=Ready pod -l name=tiller --timeout=300s

  # Install prerequisites.
  is_istio_prereq_installed="$(helm ls --all istio-init | grep istio-init)"
  if [ -z "$is_istio_prereq_installed" ]; then
    #   helm install \
    #     --wait \
    #     --name istio-init \
    #     --namespace istio-system \
    #     $istio_folder/install/kubernetes/helm/istio-init
    for i in $istio_folder/install/kubernetes/helm/istio-init/files/crd*yaml; do
      kubectl apply -f $i
    done

    # Dirty hack to let the pods install.
    sleep 15
  else
    echo "Istio prerequisites already installed."
  fi

  # Install Istio.
  echo "Installing Istio..."

  # Create namespace
  cat <<EOF | kubectl apply -f -
   apiVersion: v1
   kind: Namespace
   metadata:
     name: istio-system
     labels:
       istio-injection: disabled
EOF

  # Install Istio from a lighter template, with just pilot/gateway.
  # Based on https://knative.dev/docs/install/installing-istio/
  helm template --namespace=istio-system \
    --set prometheus.enabled=false \
    --set mixer.enabled=false \
    --set mixer.policy.enabled=false \
    --set mixer.telemetry.enabled=false \
    --set pilot.sidecar=false \
    --set pilot.resources.requests.memory=128Mi \
    --set galley.enabled=false \
    --set global.useMCP=false \
    --set security.enabled=false \
    --set global.disablePolicyChecks=true \
    --set sidecarInjectorWebhook.enabled=false \
    --set global.proxy.autoInject=disabled \
    --set global.omitSidecarInjectorConfigMap=true \
    --set gateways.istio-ingressgateway.autoscaleMin=1 \
    --set gateways.istio-ingressgateway.autoscaleMax=1 \
    --set pilot.traceSampling=100 \
    $istio_folder/install/kubernetes/helm/istio \
    >./istio-lean.yaml

  kubectl apply -f istio-lean.yaml
  rm istio-lean.yaml

  # helm install \
  #   --wait \
  #   --name istio \
  #   --namespace istio-system \
  #   $istio_folder/install/kubernetes/helm/istio \
  #   --values $istio_folder/install/kubernetes/helm/istio/values-istio-demo.yaml
}

function install_knative() {
  # Apply crd twice as workaround to:
  # https://github.com/knative/serving/issues/5722
  knative_version=0.9.0
  knative_serving_repo="https://github.com/knative/serving/releases/download"
  knative_eventing_repo="https://github.com/knative/eventing/releases/download"
  kubectl apply --wait=true --selector knative.dev/crd-install=true \
    --filename $knative_serving_repo/v$knative_version/serving.yaml \
    --filename $knative_eventing_repo/v$knative_version/release.yaml \
    --filename $knative_serving_repo/v$knative_version/monitoring.yaml ||
      \
      kubectl apply --wait=true --selector knative.dev/crd-install=true \
      --filename $knative_serving_repo/v$knative_version/serving.yaml \
      --filename $knative_eventing_repo/v$knative_version/release.yaml \
      --filename $knative_serving_repo/v$knative_version/monitoring.yaml

  kubectl apply \
    --filename $knative_serving_repo/v$knative_version/serving.yaml \
    --filename $knative_eventing_repo/v$knative_version/release.yaml \
    --filename $knative_serving_repo/v$knative_version/monitoring.yaml
}

# Create cluster if it doesn't exist.
cluster_exists=$(
  eksctl -v 0 get cluster \
    --name $cluster_name \
    --region $cluster_region
)

if [ -z "$cluster_exists" ]; then
  create_cluster
else
  echo "Cluster $cluster_name already exists."
fi

# Retrieve cluster token.
aws eks update-kubeconfig \
  --name $cluster_name \
  --region $cluster_region

# Install Nvidia drivers.
is_nvidia_installed="$(
  kubectl get daemonsets -n kube-system |
    grep nvidia-device-plugin-daemonset
)"
if [ -z "$is_nvidia_installed" ]; then
  install_nvidia_drivers
else
  echo "Nvidia drivers already installed."
fi

# Install Istio.
is_istio_installed="$(helm ls --all istio | grep 'istio ')"
if [ -z "$is_istio_installed" ]; then
  install_istio
else
  echo "Istio is already installed."
fi

# Install Knative.
is_knative_installed="$(kubectl -n knative-serving get pods | grep 'webhook')"
if [ -z "$is_knative_installed" ]; then
  install_knative
else
  echo "Knative is already installed."
fi

# Create credentials for skaffold.
# https://github.com/GoogleContainerTools/skaffold/issues/1719
is_aws_secret_setup="$(kubectl get secret | grep aws-secret)"
if [ -z "$is_aws_secret_setup" ]; then
  kubectl create secret generic aws-secret \
    --from-file /aws-credentials/credentials
else
  echo "aws-secret already configured."
fi

# Setup credential helpers in cluster for kaniko.
if [ -z "$(kubectl get secret | grep docker-kaniko-secret)" ]; then
  aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
  ecr_domain="$aws_account_id.dkr.ecr.$cluster_region.amazonaws.com"
  tmp_config_file="/tmp/config.json"
  echo "{
    \"credHelpers\": {
      \"$ecr_domain\": \"ecr-login\"
    }
  }" >$tmp_config_file
  kubectl create secret generic docker-kaniko-secret \
    --from-file $tmp_config_file
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
echo "Cluster $cluster_name is ready!"
