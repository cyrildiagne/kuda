<img src="docs/images/logo.png" width="361" height="135"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://goreportcard.com/badge/github.com/cyrildiagne/kuda?v1)](https://goreportcard.com/report/github.com/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

**Develop & deploy serverless applications on remote GPUs.**

[Kuda](https://kuda.dev) is a small util that consolidates the workflow of prototyping and deploying serverless [CUDA](https://developer.nvidia.com/cuda-zone)-based applications on [Kubernetes](http://kubernetes.io).

## Disclaimer

🧪 This is a **very early** and **experimental** work in progress:

- Don't use it in production environments.
- The API will change.

## Key Features

**Serverless GPU applications**

- Kuda uses [Knative](https://knative.dev) to consume billable GPUs only when there is traffic, and scales down to zero when there's no traffic.

**Easy to use**

- `kuda setup <provider>` : Setup a new cluster will all the requirements on the provider's managed Kubernetes, or upgrade an existing cluster.
- `kuda app deploy` : Builds & deploy an application as a serverless container.

**Language/Framework agnostic**

- Built and deployed with [Docker](https://docker.io), applications can be written in any language and use any framework.
- Applications deployed with Kuda are not required to import any specific library, keeping the code 100% portable.

**Remote development**

- The `kuda dev` command lets you spawn a remote development session with GPU inside the cluster.
- It uses [Ksync](https://github.com/vapor-ware/ksync) to synchronise your working directory with the remote session so you can code from your workstation while running the app on the remote session.

**Compatibility**

- [GCP](https://cloud.google.com) provider is already implemented.
- [AWS](https://aws.amazon.com) and [Azure](https://azure.microsoft.com) should follow soon.

## Ready?

- [Install](https://docs.kuda.dev/kuda/install)
- [Getting Started](https://docs.kuda.dev/kuda/getting_started)
- [Examples](https://github.com/cyrildiagne/kuda/tree/master/examples)
- [Reference](https://docs.kuda.dev/kuda/cli)
