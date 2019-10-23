#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
ecr_domain="$aws_account_id.dkr.ecr.$KUDA_AWS_CLUSTER_REGION.amazonaws.com"

app_name=$1
app_registry="$ecr_domain/$app_name"
app_image="$app_registry:$app_version"
namespace="default"

app_cache_name=$app_name-cache

echo $app_image

# Create Container Registry if it doesn't exists.
if [ -z "$(aws ecr describe-repositories --region $KUDA_AWS_CLUSTER_REGION | grep $app_name)" ]; then
  aws ecr create-repository \
    --repository-name $app_name \
    --region $KUDA_AWS_CLUSTER_REGION
else
  echo "Container Registry $app_registry already exists"
fi

# Create the cache registry if it doesn't exists.
if [ -z "$(aws ecr describe-repositories --region $KUDA_AWS_CLUSTER_REGION | grep $app_cache_name)" ]; then
  aws ecr create-repository \
    --repository-name $app_cache_name \
    --region $KUDA_AWS_CLUSTER_REGION
else
  echo "Container Registry $app_registry-cache already exists"
fi

# Retrieve cluster token.
aws eks update-kubeconfig \
  --name $KUDA_AWS_CLUSTER_NAME \
  --region $KUDA_AWS_CLUSTER_REGION

# Login Container Registry.
# aws ecr get-login --region $KUDA_AWS_CLUSTER_REGION --no-include-email | bash

#TODO: Build & Push image using Kaniko.

# Write Knative service config.
# cat <<EOF | kubectl apply -f -
# apiVersion: serving.knative.dev/v1alpha1
# kind: Service
# metadata:
#   name: $app_name
#   namespace: default
# spec:
#   template:
#     spec:
#       nodeSelector:
#         nvidia.com/gpu: "true"
#       tolerations:
#         - key: "nvidia.com/gpu"
#           operator: "Exists"
#           effect: "NoSchedule"
#       containers:
#         - image: $app_image
#           resources:
#             limits:
#               nvidia.com/gpu: 1
# EOF

# Write Knative service config.
echo "
apiVersion: serving.knative.dev/v1alpha1
kind: Service
metadata:
  name: $app_name
  namespace: $namespace
spec:
  template:
    spec:
      # nodeSelector:
      #   nvidia.com/gpu: true
      # tolerations:
      #   - key: nvidia.com/gpu
      #     operator: Exists
      #     effect: NoSchedule
      containers:
        - image: $app_registry
          resources:
            limits:
              nvidia.com/gpu: 1
" >.kuda-app.k8.yaml

# cat <<EOF | skaffold dev -n $namespace -f -
cat <<EOF | skaffold run -v debug -n $namespace -f -
apiVersion: skaffold/v1beta17
kind: Config
build:
  artifacts:
    - image: $app_registry
      sync:
        manual:
          - src: '**/*'
            dest: .
      kaniko:
        buildArgs:
          verbosity: debug
        buildContext:
          localDir: {}
        cache:
          repo: $app_registry-cache
        env:
          - name: AWS_REGION
            value: eu-west-1
  cluster:
    pullSecretName: aws-secret
    pullSecretMountPath: /root/.aws/
    dockerConfig: 
      secretName: docker-kaniko-secret
    namespace: $namespace
  tagPolicy:
    dateTime:
      format: "2006-01-02-15-04-05"
deploy:
  kubectl:
    manifests:
      - .kuda-app.k8.yaml
EOF
