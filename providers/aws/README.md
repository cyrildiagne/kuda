## Amazon Web Service Provider

Hacky & bare implementation with shell scripts.

# Status

| Command | Status |
| - | - |
| setup | ✔ |
| delete | ✔ |
| get | Not Started |
| app deploy | WIP |
| app delete | ✔ |
| dev start | Not Started |
| dev stop | Not Started |

# Configuration

**Prerequisites:**
- You must have subscribed to [EKS-optimize AMI with GPU support](https://aws.amazon.com/marketplace/pp/B07GRHFXGM)
- You must have an increased limit of at least 1 instance of type p2.xlarge. You can make requests [here](http://aws.amazon.com/contact-us/ec2-request)
- You must have an [aws configuration](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) on your local machine in `~/.aws/` with credentials that have authorization for:
    - cloudformation
    - ec2
    - eks
    - iam
    - api
    - ec2 autoscaling
    - ecr

# Configuration

You can override the following settings by adding them as flags of the `kuda setup` command (ex: `kuda setup aws ... --aws_...=mycluster`).

| Parameter | Default | Description |
| - | - | - |
| | | |


# Limitations
