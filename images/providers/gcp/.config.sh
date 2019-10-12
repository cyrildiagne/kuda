set -e

if [ -z "$KUDA_GCP_PROJECT_ID" ]; then
  echo "\$KUDA_GCP_PROJECT_ID is undefined."
  exit 1
fi

if [ -z "$KUDA_GCP_CREDENTIALS" ]; then
  echo "\$KUDA_GCP_CREDENTIALS is undefined."
  exit 1
fi

echo
echo -e "\e[1m \e[34mKuda GCP provider \e[36mv$(cat /kuda_cmd/VERSION) \e[0m"
echo

export KUDA_GCP_CREDENTIALS=/secret/$(basename $KUDA_GCP_CREDENTIALS)

# Set default config.
export KUDA_GCP_CLUSTER_NAME="${KUDA_GCP_CLUSTER_NAME:-kuda}"
export KUDA_GCP_COMPUTE_ZONE="${KUDA_GCP_COMPUTE_ZONE:-us-central1-a}"
export KUDA_GCP_MACHINE_TYPE="${KUDA_GCP_MACHINE_TYPE:-n1-standard-2}"

export KUDA_DEFAULT_POOL_NUM_NODES="${KUDA_DEFAULT_POOL_NUM_NODES:-1}"
export KUDA_DEFAULT_GPU="${KUDA_DEFAULT_GPU:-k80}"
export KUDA_DEFAULT_USE_PREEMPTIBLE="${KUDA_DEFAULT_USE_PREEMPTIBLE:-false}"

export KUDA_DEV_APP_NAME="${KUDA_DEV_APP_NAME:-kuda-dev}"
export KUDA_DEV_SYNC_PATH="${KUDA_DEV_SYNC_PATH:-/app_home}"

# Disable prompts.
gcloud config set survey/disable_prompts true

# Apply to gcloud config.
gcloud config set project $KUDA_GCP_PROJECT_ID
gcloud config set compute/zone $KUDA_GCP_COMPUTE_ZONE

# Enable service account from the credential json.
credentials_file=$(basename "$KUDA_GCP_CREDENTIALS")
gcloud auth activate-service-account --key-file /secret/$credentials_file
