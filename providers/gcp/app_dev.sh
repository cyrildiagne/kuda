#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

app_name=$1
app_image="gcr.io/$KUDA_GCP_PROJECT_ID/$app_name"
namespace="default"

# Retrieve cluster token.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# Login Docker (needed for skaffold's sync).
gcloud auth configure-docker --quiet

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
      containers:
        - image: $app_image
          command: ["python3"]
          args: ["app.py"]
          # resources:
          #   limits:
          #     nvidia.com/gpu: 1
          volumeMounts:
            - name: secret
              readOnly: true
              mountPath: "/secret"
          env:
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /secret/$(basename $KUDA_GCP_CREDENTIALS)
      volumes:
        - name: secret
          secret:
            secretName: $(basename $KUDA_GCP_CREDENTIALS)
" > .kuda-app.k8.yaml

# Cloud Build has a generous free tier is easy enough to use with Skaffold
# So we use it rather than Kaniko.
export GOOGLE_APPLICATION_CREDENTIALS=$KUDA_GCP_CREDENTIALS
cat <<EOF | skaffold dev -n $namespace -f -
apiVersion: skaffold/v1beta16
kind: Config
build:
  googleCloudBuild:
    projectId: $KUDA_GCP_PROJECT_ID
  artifacts:
    - image: $app_image
      sync:
        manual:
          - src: '**/*'
            dest: './'
  tagPolicy:
    dateTime:
      format: "2006-01-02_15-04-05"
deploy:
  kubectl:
    manifests:
      - .kuda-app.k8.yaml
EOF

rm .kuda-app.k8.yaml