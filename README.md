<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

🧪 **Status:** early and experimental WIP

---

## Rapidly develop and deploy serverless APIs that need GPUs

Kuda deploys APIs as serverless containers using [Knative](https://knative.dev)
which means that you can use any language and any framework, and there is no library to import in your code.
All you need is a Dockerfile.

## A simple interface for the full serverless API development cycle

- `kuda dev` Deploy the API on remote GPUs in dev mode (with file sync & live reload)
- `kuda deploy` Deploy the API in production mode.
  It will be automatically scaled down to zero when there is no traffic,
  and back up when there are new requests.

## Run your APIs anywhere Kubernetes is running

<!-- - [gpu.sh](#) - The best way to get started quickly on a cost-effective, fully-managed GPU cluster. -->

- [GKE](#) - Installation guide for [running Kuda on GCP](/docs/install_on_gcp.md).

## Use the tools you know

Here's a simple example that prints the result of `nvidia-smi` using [Flask](http://flask.palletsprojects.com):

- `main.py`

```python
import os
import flask

app = flask.Flask(__name__)

@app.route('/')
def hello():
    return 'Hello GPU!\n\n' + os.popen('nvidia-smi').read()

if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=80)
```

- `Dockerfile`

```Dockerfile
FROM nvidia/cuda:10.1-base

RUN apt-get update && apt-get install -y --no-install-recommends \
  python3 python3-pip \
  && \
  apt-get clean && \
  apt-get autoremove && \
  rm -rf /var/lib/apt/lists/*
RUN pip3 install setuptools Flask gunicorn

WORKDIR /app

COPY main.py ./main.py

CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 main:app
```

- `kuda.yaml`

```yaml
kudaManifestVersion: v1alpha1

name: hello-gpu

deploy:
  dockerfile: ./Dockerfile
```

This `kuda deploy` build and deploy the API which you can call via HTTP,
for instance with [cURL](https://curl.haxx.se/):

```bash
$ curl https://hello-gpu.kuda.yourdomain.com

Hello GPU!

+-----------------------------------------------------------------------------+
| NVIDIA-SMI 418.67       Driver Version: 418.67       CUDA Version: 10.1     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  Tesla K80           Off  | 00000000:00:04.0 Off |                    0 |
| N/A   37C    P8    27W / 149W |      0MiB / 11441MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+

```

## Get Started

- [Install](docs/install_cli.md)
- [Getting Started](docs/getting_started.md)
- [CLI Reference](docs/cli.md)
