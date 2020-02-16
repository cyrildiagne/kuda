<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

**Status:** Experimental 🧪

## Easily deploy GPU models as serverless APIs

- Deploy a template

```bash
$ kuda deploy -f https://github.com/cyrildiagne/kuda/releases/download/v0.4.0-preview/example-hello-gpu.yaml
```

- Call it!

```bash
$ curl https://hello-gpu.default.$your_domain
```

<!-- ```bash
$ curl -H 'x-api-key: $your_key' https://hello-gpu.default.$your_domain
``` -->

```
Hello GPU!

+-----------------------------------------------------------------------------+
| NVIDIA-SMI 418.67       Driver Version: 418.67       CUDA Version: 10.1     |
|-------------------------------+----------------------+----------------------+
| GPU  Name        Persistence-M| Bus-Id        Disp.A | Volatile Uncorr. ECC |
| Fan  Temp  Perf  Pwr:Usage/Cap|         Memory-Usage | GPU-Util  Compute M. |
|===============================+======================+======================|
|   0  Tesla T4            Off  | 00000000:00:04.0 Off |                    0 |
| N/A   37C    P8    10W /  70W |      0MiB / 15079MiB |      0%      Default |
+-------------------------------+----------------------+----------------------+

+-----------------------------------------------------------------------------+
| Processes:                                                       GPU Memory |
|  GPU       PID   Type   Process name                             Usage      |
|=============================================================================|
|  No running processes found                                                 |
+-----------------------------------------------------------------------------+

```

<!-- ## Add GPU models to a webapp easily

- Deploy a template from github

```bash
kuda deploy -f cyrildiagne/gpt2-small-http
```

- Call your deployed API

```bash
$ curl \
    -H 'x-api-key: $your_key' \
    -F 'input=Kuda is' \
    https://gpt2.<your-namespace>.example.com/generate
```

```json
{
  "query": "Kuda is",
  "generated": "a tool that...etc."
}
```

Checkout the full list of templates available in [the registry](#).
-->

## Serverless GPU inference

Kuda builds on [Knative](#) to allocate cloud GPUs only when there is traffic
to your app.

This is ideal when you want to share ML projects online without keeping
expensive GPUs allocated all the time.

It tries to reduce cold starts time (gpu nodes allocation and service instanciation)
as much possible and to tries manage cooldown times intelligently.

## Turn any model into a serverless API

Kuda deploys APIs as a docker containers, so you can use any language, any
framework, and there is no library to import in your code.

All you need is a Dockerfile.

Here's a minimal example that just prints the result of `nvidia-smi` using
[Flask](http://flask.palletsprojects.com):

- `main.py`

```python
import os
import flask

app = flask.Flask(__name__)

@app.route('/')
def hello():
    return 'Hello GPU:\n' + os.popen('nvidia-smi').read()
```

- `Dockerfile`

```Dockerfile
FROM nvidia/cuda:10.1-base

RUN apt-get install -y python3 python3-pip

RUN pip3 install setuptools Flask gunicorn

COPY main.py ./main.py

CMD exec gunicorn --bind :80 --workers 1 --threads 8 main:app
```

- `kuda.yaml`

```yaml
name: hello-gpu
deploy:
  dockerfile: ./Dockerfile
```

Running `kuda deploy` in this example would build and deploy the API to a url
such as `https://hello-gpu.my-namespace.example.com`.

Checkout the full example with annotations in
[examples/hello-gpu-flask](examples/hello-gpu-flask).

## Features

- Provision GPUs & scale based on traffic (from zero to N)
- Interactive development on cloud GPUs from any workstation
- Protect & control access to your APIs using API Keys
- HTTPS with TLS termination & automatic certificate management

## Get Started

- [Install](docs/install_cli.md)
- [Getting Started](docs/getting_started.md)
- [CLI Reference](docs/cli.md)
