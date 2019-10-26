<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://goreportcard.com/badge/github.com/cyrildiagne/kuda?v1)](https://goreportcard.com/report/github.com/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

**Develop & deploy serverless applications on remote GPUs.**

[Kuda](https://kuda.dev) helps prototyping and deploying serverless applications that need GPUs and [CUDA](https://developer.nvidia.com/cuda-zone) on [Kubernetes](http://kubernetes.io).

It is based on [Knative](https://knative.dev), [Skaffold](https://skaffold.dev) and [Kaniko](https://github.com/GoogleContainerTools/kaniko), and supports the major cloud providers.

## Disclaimer

🧪 This is a **very early** and **experimental** work in progress:

- Most things won't work out of the box.
- It might break things in the cluster. Keep it away from production resources :)

## Key Features

**Serverless GPU applications**

- Kuda uses [Knative](https://knative.dev) to consume GPUs only when there is traffic, and scales down to zero when there's no traffic.

**Easy to use**

- `kuda setup <provider>` : Setup a new cluster will all the requirements on the provider's managed Kubernetes, or upgrade an existing cluster.
- `kuda app dev` : Deploys an application and watches your local folder so that the app reloads automatically on the cluster when you make local changes.
- `kuda app deploy` : Deploy the application as a serverless container.

**Language/Framework agnostic**

- Built and deployed with [Docker](https://docker.io), applications can be written in any language and use any framework.
- Applications deployed with Kuda are not required to import any specific library.

**Cloud provider Compatibility**

| Provider | Status                          |
| -------- | ------------------------------- |
| GCP      | [In progress...](providers/gcp) |
| AWS      | [In progress...](providers/aws) |
| Azure    | Not started                     |
| NGC      | Not started                     |

## Ready?

- [Install](docs/kuda/install.md)
- [Getting Started](docs/kuda/getting_started.md)
- [Examples](https://github.com/cyrildiagne/kuda-apps)
- [CLI Reference](docs/kuda/cli.md)
