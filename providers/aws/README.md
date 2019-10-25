## Amazon Web Service Provider

Hacky & bare implementation with shell scripts.
It uses [ECR](https://aws.amazon.com/ecr) to privately store the app images.

# Status

| Command      |  Status |
| ------------ | ------- |
| `setup`      | ✔       |
| `delete`     | ✔       |
| `app dev`    | ✔       |
| `app deploy` | ✔       |
| `app delete` | ✔       |

Functionalities:

| Functionality       |  Status     |
| ------------------- | ----------- |
| Dev file sync       | ✔           |
| GPU node autoscaler | WIP         |
| Https               | Not started |
| Dns                 | Not started |
| Monitoring          | Not started |

Status Notes & current issues:

- https://github.com/cyrildiagne/kuda/issues/11

# Prerequisites

- You must have subscribed to [EKS-optimize AMI with GPU support](https://aws.amazon.com/marketplace/pp/B07GRHFXGM)
- You must have an increased limit of at least 1 instance of type p2.xlarge. You can make requests [here](http://aws.amazon.com/contact-us/ec2-request)
- You must have an [aws configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) on your local machine in `~/.aws/` with credentials that have authorization for:
  - cloudformation
  - ec2
  - ec2 autoscaling
  - eks
  - iam
  - api
  - ecr
  - elb

<!-- # Configuration

You can override the following settings by adding them as flags of the `kuda setup` command (ex: `kuda setup aws ... --aws_...=mycluster`).

| Parameter | Default | Description |
| - | - | - |
| | | |


# Limitations -->
