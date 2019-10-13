FROM alpine:3.10

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
RUN curl --silent \
  --location "https://github.com/weaveworks/eksctl/releases/download/latest_release/eksctl_$(uname -s)_amd64.tar.gz" \
  | tar xz -C /tmp
RUN mv /tmp/eksctl /usr/local/bin

# Install kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

# Install Helm
RUN curl -L https://git.io/get_helm.sh | bash

# Install Ksync.
RUN curl https://vapor-ware.github.io/gimme-that/gimme.sh | bash

# Download Istio CRDs
RUN curl -L https://git.io/getLatestIstio | sh -

# Copy the provider's commands.
COPY . /kuda_cmd
ENV KUDA_CMD_DIR /kuda_cmd
RUN chmod +x /kuda_cmd/*.sh
RUN ln -s /kuda_cmd/app_deploy.sh /usr/local/bin/kuda_app_deploy && \
  ln -s /kuda_cmd/app_delete.sh /usr/local/bin/kuda_app_delete && \
  ln -s /kuda_cmd/setup.sh /usr/local/bin/kuda_setup && \
  ln -s /kuda_cmd/delete.sh /usr/local/bin/kuda_delete && \
  ln -s /kuda_cmd/dev_start.sh /usr/local/bin/kuda_dev_start && \
  ln -s /kuda_cmd/dev_stop.sh /usr/local/bin/kuda_dev_stop && \
  ln -s /kuda_cmd/get.sh /usr/local/bin/kuda_get

# Go to the app home.
WORKDIR /app_home

ENTRYPOINT ["/bin/bash"]