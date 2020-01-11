<img src="docs/images/logo.png" width="241" height="90"/>

[![](https://circleci.com/gh/cyrildiagne/kuda/tree/master.svg?style=shield&circle-token=b14f5838ae2acabe21a8255070507f7e36ba510b)](https://circleci.com/gh/cyrildiagne/kuda)
[![](https://img.shields.io/github/v/release/cyrildiagne/kuda?include_prereleases)](https://github.com/cyrildiagne/kuda/releases)

**Status:** 🧪Experimental

## Turn any model into a serverless API

Easily turn any model into a serverless API that will consume cloud GPUs only
when it's being called.

Kuda deploys your API as a docker container, so you can use any language, any
framework, and there is no library to import in your code.

All you need is a Dockerfile.

## Easy to use

- `kuda init` Initializes your local & remote configurations
- `kuda dev` Deploys the API in dev mode (with file sync & live reload)
- `kuda deploy` Deploys the API in production mode
- `kuda publish` Publishes the API template to the registry

## Features

- Provision GPUs & scale based on traffic (from zero to N)
- Interactive development on cloud GPUs from any workstation
- Protect & control access to your APIs using API Keys
- HTTPS with TLS termination & automatic certificate management

## Use the frameworks you know

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
such as `https://hello-gpu.default.yourdomain.com` which you can easily call:

```
$ curl https://hello-gpu.default.yourdomain.com

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

Checkout the full example with annotations in
[examples/hello-gpu-flask](examples/hello-gpu-flask).

## Get Started

- [Install](docs/install_cli.md)
- [Getting Started](docs/getting_started.md)
- [CLI Reference](docs/cli.md)
