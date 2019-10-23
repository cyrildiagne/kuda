#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

app_name=$(echo $1 | cut -f1 -d':')
app_version=$(echo $1 | cut -f2 -d':')
# app_image="gcr.io/$KUDA_GCP_PROJECT_ID/$app_name:$app_version"
app_image="gcr.io/$KUDA_GCP_PROJECT_ID/$app_name"
namespace="default"
echo $app_image

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

export GOOGLE_APPLICATION_CREDENTIALS=$KUDA_GCP_CREDENTIALS

curl -fsSL "https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v1.5.0/docker-credential-gcr_linux_amd64-1.5.0.tar.gz" | tar xz --to-stdout ./docker-credential-gcr > /usr/bin/docker-credential-gcr && chmod +x /usr/bin/docker-credential-gcr
docker-credential-gcr configure-docker

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
          resources:
            limits:
              nvidia.com/gpu: 1
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

cat <<EOF | skaffold dev -v debug --cache-artifacts=false -n $namespace -f -
apiVersion: skaffold/v1beta16
kind: Config
build:
  googleCloudBuild:
    projectId: $KUDA_GCP_PROJECT_ID
  artifacts:
    - image: $app_image
      sync:
        manual:
          - src: './**/*'
            dest: .
  tagPolicy:
    dateTime:
      format: "2006-01-02_15-04-05.999_MST"
      timezone: "Local"
deploy:
  kubectl:
    manifests:
      - .kuda-app.k8.yaml
EOF

rm .kuda-app.k8.yaml