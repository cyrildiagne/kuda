#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# Config.
cluster_name=$KUDA_AWS_CLUSTER_NAME
aws_region=$KUDA_AWS_CLUSTER_REGION
aws_account_id="$(aws sts get-caller-identity | jq -r .Account)"
ecr_domain="$aws_account_id.dkr.ecr.$aws_region.amazonaws.com"
app_name=$(echo $1 | cut -f1 -d':')
app_version=$(echo $1 | cut -f2 -d':')
app_registry="$ecr_domain/$app_name"
namespace="default"
app_cache_name=$app_name-cache

# Create container registries if they doesn't exists.
repos="$(aws ecr describe-repositories --region $aws_region)"
if [ -z "$(echo $repos | grep $app_name)" ]; then
  aws ecr create-repository \
    --repository-name $app_name \
    --region $aws_region
else
  echo "Container Registry $app_registry already exists"
fi
# Cache registry.
if [ -z "$(echo $repos | grep $app_cache_name)" ]; then
  aws ecr create-repository \
    --repository-name $app_cache_name \
    --region $aws_region
else
  echo "Container Registry $app_registry-cache already exists"
fi


# Retrieve cluster token.
aws eks update-kubeconfig \
  --name $cluster_name \
  --region $aws_region

# Login Container Registry.
aws ecr get-login --region $aws_region --no-include-email | bash

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
      #   nvidia.com/gpu: 'true'
      # tolerations:
      #   - key: 'nvidia.com/gpu'
      #     operator: 'Exists'
      #     effect: 'NoSchedule'
      containers:
        - image: $app_registry
          resources:
            limits:
              nvidia.com/gpu: 1
          env:
            - name: KUDA_DEV
              value: 'true'
" >.kuda-app.k8.yaml

# Send version as env variable for Skaffold to use.
export APP_VERSION=$app_version
cat <<EOF | skaffold run -n $namespace -f -
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
    envTemplate:
      template: "{{.IMAGE_NAME}}:{{.APP_VERSION}}"
deploy:
  kubectl:
    manifests:
      - .kuda-app.k8.yaml
EOF
rm .kuda-app.k8.yaml
