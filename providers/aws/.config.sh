set -e

echo
echo -e "\e[1m \e[34mKuda AWS provider \e[36mv$(cat /kuda_cmd/VERSION) \e[0m"
echo

# Set custom path to AWS config file
export AWS_CONFIG_FILE='/aws-credentials/config'
export AWS_SHARED_CREDENTIALS_FILE='/aws-credentials/credentials'

# Set KUDA envs.
export KUDA_AWS_CLUSTER_NAME="${KUDA_CLUSTER_NAME:-kuda}"
export KUDA_AWS_CLUSTER_REGION="${KUDA_CLUSTER_REGION:-eu-west-1}"