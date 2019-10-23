FROM alpine:3.10

ARG CLOUD_SDK_VERSION=260.0.0
ARG SKAFFOLD_VERSION=0.40.0

# Install Google Cloud SDK and Docker.
# Docker is needed to enable Skaffold dev.
ENV CLOUD_SDK_VERSION=$CLOUD_SDK_VERSION
ENV PATH /google-cloud-sdk/bin:$PATH
RUN apk --no-cache add \
  docker \
  bash \
  curl \
  python \
  py-crcmod \
  libc6-compat \
  openssh-client \
  git \
  gnupg \
  && curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
  tar xzf google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
  rm google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz && \
  gcloud config set core/disable_usage_reporting true && \
  gcloud config set component_manager/disable_update_check true && \
  gcloud config set metrics/environment github_docker_image && \
  gcloud --version
VOLUME ["/root/.config"]

# Install Kubectl.
RUN gcloud components install beta kubectl

# Install Skaffold.
RUN curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/v${SKAFFOLD_VERSION}/skaffold-linux-amd64 && \
  chmod +x skaffold && \
  mv skaffold /usr/local/bin

# Copy the provider's commands.
COPY . /kuda_cmd
ENV KUDA_CMD_DIR /kuda_cmd
RUN chmod +x /kuda_cmd/*.sh
RUN ln -s /kuda_cmd/app_deploy.sh /usr/local/bin/kuda_app_deploy && \
  ln -s /kuda_cmd/app_dev.sh /usr/local/bin/kuda_app_dev && \
  ln -s /kuda_cmd/app_delete.sh /usr/local/bin/kuda_app_delete && \
  ln -s /kuda_cmd/setup.sh /usr/local/bin/kuda_setup && \
  ln -s /kuda_cmd/delete.sh /usr/local/bin/kuda_delete

# Go to the app home.
WORKDIR /app_home

ENTRYPOINT ["/bin/bash"]