FROM docker:17.12.0-ce as static-docker-source

FROM golang:1.13.5 as builder

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    mv ./kubectl /tmp/kubectl

RUN curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64 && \
    mv ./skaffold /tmp/skaffold

COPY go.mod go.sum /go/src/github.com/cyrildiagne/kuda/
WORKDIR /go/src/github.com/cyrildiagne/kuda
RUN go mod download

COPY pkg/config /go/src/github.com/cyrildiagne/kuda/pkg/config
COPY pkg/manifest /go/src/github.com/cyrildiagne/kuda/pkg/manifest
COPY pkg/utils /go/src/github.com/cyrildiagne/kuda/pkg/utils
COPY pkg/deployer /go/src/github.com/cyrildiagne/kuda/pkg/deployer
WORKDIR /go/src/github.com/cyrildiagne/kuda/pkg/deployer
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o deployer .

#

FROM alpine:3.11

ARG CLOUD_SDK_VERSION=274.0.1
ENV CLOUD_SDK_VERSION=$CLOUD_SDK_VERSION

ENV PATH /google-cloud-sdk/bin:$PATH
COPY --from=static-docker-source /usr/local/bin/docker /usr/local/bin/docker
RUN apk --no-cache add \
        ca-certificates \
        curl \
        python \
        py-crcmod \
        bash \
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

COPY --from=builder /tmp/kubectl /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl

COPY --from=builder /tmp/skaffold /usr/local/bin/skaffold
RUN chmod +x /usr/local/bin/skaffold

COPY --from=builder /go/src/github.com/cyrildiagne/kuda/pkg/deployer/deployer /root/
RUN chmod +x /root/deployer

# Launch the app on port 80.
ENV PORT 80

ENTRYPOINT ["/root/deployer"]