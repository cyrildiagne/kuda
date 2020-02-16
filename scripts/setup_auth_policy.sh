#!/bin/bash

# Make sure you've configured the auth module before running this script
# otherwise you won't be able to the access services in the kuda namespace.

set -e

source "$(dirname $BASH_SOURCE)/utils.sh"

# If using Firebase / Cloud Identity as authentication provider.
FIREBASE_JWT_ISSUER="https://securetoken.google.com/$KUDA_GCP_PROJECT"
FIREBASE_JWT_URI="https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

export KUDA_JWT_ISSUER="${KUDA_JWT_ISSUER:-$FIREBASE_JWT_ISSUER}"
export KUDA_JWT_URI="${KUDA_JWT_URI:-$FIREBASE_JWT_URI}"

assert_set KUDA_GCP_PROJECT $KUDA_GCP_PROJECT
assert_set KUDA_JWT_ISSUER $KUDA_JWT_ISSUER
assert_set KUDA_JWT_URI $KUDA_JWT_URI

function setup_deployer_auth_policy() {
  echo "Adding istio authentication policy for deployer..."

  kubectl apply -f - <<EOF
apiVersion: authentication.istio.io/v1alpha1
kind: Policy
metadata:
  name: api-origin-auth
  namespace: kuda
spec:
  targets:
  - name: api
    ports:
    - number: 80
    - number: 443
  origins:
  - jwt:
      issuer: $KUDA_JWT_ISSUER
      audiences:
      - "$KUDA_GCP_PROJECT"
      jwksUri: "$KUDA_JWT_URI"
      triggerRules:
      - excluded_paths:
        - prefix: /metrics
        - prefix: /healthz
  principalBinding: USE_ORIGIN
EOF
}

setup_deployer_auth_policy
