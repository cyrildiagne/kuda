set -e

echo
echo -e "\e[1m \e[34mKuda GCP provider \e[36mv$(cat /kuda_cmd/VERSION) \e[0m"
echo

# Make sure gcp_project_id is set.
if [ -z "$KUDA_GCP_PROJECT_ID" ]; then
  echo "\$KUDA_GCP_PROJECT_ID is undefined."
  exit 1
fi

# Make sure gcp_credentials is set.
if [ -z "$KUDA_GCP_CREDENTIALS" ]; then
  echo "\$KUDA_GCP_CREDENTIALS is undefined."
  exit 1
fi

# Setup credential path to mounted volume in the docker image.
export KUDA_GCP_CREDENTIALS=/secret/$(basename $KUDA_GCP_CREDENTIALS)

# Set default Kuda Dev config.
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
