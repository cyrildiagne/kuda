FROM alpine:3.10

# Set dependencies versions.
ARG EKSCTL_VERSION=0.7.0
ARG KUBECTL_VERSION=1.16.2
ARG SKAFFOLD_VERSION=0.40.0
ARG ISTIO_VERSION=1.3.3
ARG HELM_VERSION=2.15.0

# Install base apps.
RUN apk add --update-cache --no-cache \
  bash \
  build-base \
  curl \
  python \
  python-dev \
  py-pip \
  openssl \
  jq

# Install AWS CLI
RUN pip install awscli

# Install eksctl
ARG EKSCTL_HOST="https://github.com/weaveworks/eksctl/releases/download"
RUN curl --silent \
  --location "${EKSCTL_HOST}/${EKSCTL_VERSION}/eksctl_$(uname -s)_amd64.tar.gz" \
  | tar xz -C /tmp
RUN mv /tmp/eksctl /usr/local/bin

# Install kubectl
ARG KUBECTL_HOST="https://storage.googleapis.com/kubernetes-release/release"
RUN curl -LO "${KUBECTL_HOST}/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" && \
  chmod +x ./kubectl && \
  mv ./kubectl /usr/local/bin/kubectl

# Install Helm
RUN curl -L https://git.io/get_helm.sh | DESIRED_VERSION="v${HELM_VERSION}" bash

# Download Istio CRDs
RUN curl -L https://git.io/getLatestIstio | ISTIO_VERSION="${ISTIO_VERSION}" sh -

# Install Skaffold.
ARG SKAFFOLD_HOST="https://storage.googleapis.com/skaffold/releases"
# RUN curl -Lo skaffold "${SKAFFOLD_HOST}/v${SKAFFOLD_VERSION}/skaffold-linux-amd64" && \
RUN curl -Lo skaffold "https://storage.googleapis.com/artifacts.gpu-sh.appspot.com/releases/latest/skaffold-linux-amd64" && \
  chmod +x skaffold && \
  mv skaffold /usr/local/bin

# Copy the provider's commands.
COPY . /kuda_cmd
ENV KUDA_CMD_DIR /kuda_cmd
RUN chmod +x /kuda_cmd/*.sh
RUN ln -s /kuda_cmd/app_deploy.sh /usr/local/bin/kuda_app_deploy && \
  ln -s /kuda_cmd/app_delete.sh /usr/local/bin/kuda_app_delete && \
  ln -s /kuda_cmd/app_dev.sh /usr/local/bin/kuda_app_dev && \
  ln -s /kuda_cmd/setup.sh /usr/local/bin/kuda_setup && \
  ln -s /kuda_cmd/delete.sh /usr/local/bin/kuda_delete && \
  ln -s /kuda_cmd/get.sh /usr/local/bin/kuda_get

# Go to the app home.
WORKDIR /app_home

ENTRYPOINT ["/bin/bash"]