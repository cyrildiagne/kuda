# Remote Development

**⚠️ Remote development and this guide are still WIP. Following this guide probably won't work for now.**

This guide will walk you through the process of developping remotely on the Kubernetes cluster.

Make sure you have a cluster running with Kuda's dependencies.

## 1 - Introduction

Remote dev sessions work like a virtual machine running in you kubernetes cluster that will use [Ksync](https://github.com/vapor-ware/ksync/) to synchronize the local folder on your workstation. So you can code from any machine with your favorite IDE while running the workloads on powerful remote GPUs.

Developping on remote sessions offers many other advantages such as:

- Elastic resources - Scale up and down the hardware for the current task.
- Datacenter-fast internet - Download large datasets _much_ faster.
- Contained environment per project - No more conflict between librairies, CUDA or python version..etc.

## 2 - Start a Remote Dev Session

First clone the kuda-apps repository and navigate to the `hello-gpu` example.
```bash
git clone https://github.com/cyrildiagne/kuda-apps
cd hello-gpu
```

Start a remote dev session that will be provisioned on your cluster.

```bash
kuda dev start gcr.io/deeplearning-platform-release/base-cu100
```

`gcr.io/deeplearning-platform-release/base-cu100` Is the docker image to use as base. This image is convenient if you're using Kuda for deep learning since it packages most of the softwares needed in the deeplearning development cycle. It also allows you to specify which version of CUDA and CuDNN you need.

This command will start the remote session and synchronize the CWD \(current working directory\) with the remote instance.

Upon started, it will also print the cluster's IP address / port to use later on. Make not of that as we'll refer to it later as `<your-dev-session-external-ip:port>`

List of recommended `base-image`:

- all images from [nvidia/cuda](https://hub.docker.com/r/nvidia/cuda/). These images are fairly lightweight but python must be installed manually.
- gcloud's [Deep Learning containers](https://cloud.google.com/ai-platform/deep-learning-containers/docs/choosing-container)

→ For more information on the `kuda dev start` command, check the [reference](https://docs.kuda.dev/kuda/cli#dev).

## • Retrieve & initialize an example application

The next command should be ran in the remote shell that's started.

Install the example dependencies. Because the remote dev session is short lived, you don't need to create a virtualenv:

```bash
root@kuda-dev: pip install -r requirements.txt
```

## • Run and Test the example

Then start the example in dev mode. It will reload automatically when you make changes from your local machine:

```bash
root@kuda-dev: python app.py
```

Open `http://<your-dev-session-external-ip:port>/101` in a web browser. You should see the result of the `nvidia-smi` command which gives some details about the GPU.

You can make a change to the code of the example on your local workstation and the application will restart automatically after the code has synched.

## 3 - Cleanup

### • Stop the remote dev session

From the remote shell that's opened after you've run `kuda dev start`, simply enter:

```bash
root@kuda-dev: exit
```

This will automatically call `kuda dev stop` and shut down the remote session.
