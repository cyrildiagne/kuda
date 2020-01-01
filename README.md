<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://goreportcard.com/badge/github.com/cyrildiagne/kuda?v1)](https://goreportcard.com/report/github.com/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

**Develop, deploy and manage serverless APIs that need GPUs on Kubernetes.**

Kuda is based on [Knative](https://knative.dev) and
[Skaffold](https://skaffold.dev) and provides a simple interface for the full API development cycle:

- `kuda init <name>` Initialize the API configuration files
- `kuda dev` Deploy the API in dev mode (with file sync & live reload)
- `kuda deploy` Deploy the API in production mode

## Disclaimer

🧪 It's an **early** and **experimental** work in progress.

## Getting Started

- [Install](docs/install.md)
- [Getting Started](docs/getting_started.md)
- [CLI Reference](docs/cli.md)