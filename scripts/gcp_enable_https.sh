#!/bin/bash

# You must have a real domain name (not xip.io auto-domain) to enable HTTPS.
# This scripts enables HTTPS using CloudDNS, Let's Encrypt and cert-manager.
# Adapt the ClusterIssuer manifest if you are using a different DNS.

set -e

red="\033[31m"
reset="\033[0m"

function print_help_and_exit() {
  echo "
This script requires 4 environment variables to be set:
- KUDA_GCP_PROJECT: The gcloud project ID for CloudDNS.
- KUDA_DOMAIN: Your domain name.
- KUDA_NAMESPACE: Your Kuda namespace.
- KUDA_LETSENCRYPT_EMAIL: The admin email for Let's Encrypt.

Example usage:
  export KUDA_GCP_PROJECT=your-gcp-project
  export KUDA_DOMAIN=example.com
  export KUDA_NAMESPACE=default
  export KUDA_LETSENCRYPT_EMAIL=you@example.com
  sh scripts/gcp_enable_https.sh
"
  exit 1
}

function assert_set() {
  var_name=$1
  var_value=$2
  if [ -z "$var_value" ]; then
    printf "${red}ERROR:${reset} Missing required env variable $var_name\n"
    print_help_and_exit
  else
    echo "Using $var_name: $var_value"
  fi
}

assert_set KUDA_GCP_PROJECT $KUDA_GCP_PROJECT
assert_set KUDA_DOMAIN $KUDA_DOMAIN
assert_set KUDA_NAMESPACE $KUDA_NAMESPACE
assert_set KUDA_LETSENCRYPT_EMAIL $KUDA_LETSENCRYPT_EMAIL

# Cert Manager Version
CERT_MANAGER_VERSION=0.12.0
# Knative Version
KNATIVE_VERSION=0.11.0
# Name of the service account you want to create.
CLOUD_DNS_SA=certm-cdns-admin-$(date +%s)
# Fully-qualified service account name also has project-id information.
CLOUD_DNS_SA_EMAIL=$CLOUD_DNS_SA@$KUDA_GCP_PROJECT.iam.gserviceaccount.com

# Install certmanager CRDs & resources.
function install_cert_manager() {
  certmanager_repo="https://github.com/jetstack/cert-manager"
  kubectl create namespace cert-manager
  kubectl apply --validate=false \
    -f $certmanager_repo/releases/download/v$CERT_MANAGER_VERSION/cert-manager.yaml
}

# Install knative certmanager component.
function install_knative_cert_serving() {
  knative_repo="https://github.com/knative/serving"
  kubectl apply \
    -f $knative_repo/releases/download/v$KNATIVE_VERSION/serving-cert-manager.yaml
}

function create_service_account() {
  if gcloud --project $KUDA_GCP_PROJECT iam service-accounts list | grep $CLOUD_DNS_SA; then
    echo "service account already exists."
  else
    gcloud --project $KUDA_GCP_PROJECT iam service-accounts \
      create $CLOUD_DNS_SA \
      --display-name "Service Account to support ACME DNS-01 challenge."

    # Bind the role dns.admin to this service account, so it can be used to support
    # the ACME DNS01 challenge.
    gcloud projects add-iam-policy-binding $KUDA_GCP_PROJECT \
      --member serviceAccount:$CLOUD_DNS_SA_EMAIL \
      --role roles/dns.admin
  fi
}

function create_service_account_secret() {
  # Make a temporary directory to store key
  KEY_DIRECTORY=$(mktemp -d)

  # Download the secret key file for your service account.
  gcloud iam service-accounts keys create $KEY_DIRECTORY/cloud-dns-key.json \
    --iam-account=$CLOUD_DNS_SA_EMAIL

  # Upload that as a secret in your Kubernetes cluster.
  kubectl create secret -n cert-manager generic cloud-dns-key \
    --from-file=key.json=$KEY_DIRECTORY/cloud-dns-key.json

  # Delete the local secret
  rm -rf $KEY_DIRECTORY
}

function create_cluster_issuer() {
  kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: ClusterIssuer
metadata:
  name: letsencrypt-issuer
  namespace: cert-manager
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: $KUDA_LETSENCRYPT_EMAIL
    privateKeySecretRef:
      name: letsencrypt-issuer
    solvers:
    - selector: {}
      dns01:
        clouddns:
          project: $KUDA_GCP_PROJECT
          serviceAccountSecretRef:
            name: cloud-dns-key
            key: key.json
EOF
}

function create_certificate() {
  kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: $KUDA_GCP_PROJECT
  # Istio certs secret lives in the istio-system namespace, and
  # a cert-manager Certificate is namespace-scoped.
  namespace: istio-system
spec:
  # Reference to the Istio default cert secret.
  secretName: istio-ingressgateway-certs
  # The certificate common name, use one from your domains.
  commonName: "*.kuda.$KUDA_DOMAIN"
  dnsNames:
    # Since certificate wildcards only allow one level, we will
    # need to one for every namespace that Knative is used in.
    # We don't need to use wildcard here, fully-qualified domains
    # will work fine too.
    - "*.kuda.$KUDA_DOMAIN"
    - "*.$KUDA_NAMESPACE.$KUDA_DOMAIN"
  # Reference to the ClusterIssuer we created in the previous step.
  issuerRef:
    kind: ClusterIssuer
    name: letsencrypt-issuer
EOF
}

function update_knative_ingress_gateway() {
  kubectl apply -f - <<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: knative-ingress-gateway
  namespace: knative-serving
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
  - port:
      number: 443
      name: https
      protocol: HTTPS
    hosts:
    - "*"
    tls:
      mode: SIMPLE
      privateKey: /etc/istio/ingressgateway-certs/tls.key
      serverCertificate: /etc/istio/ingressgateway-certs/tls.crt
EOF
}

echo "Installing Cert Manager"
install_cert_manager

echo "Install Knative Serving cert component.."
install_knative_cert_serving

echo "Creating service account $CLOUD_DNS_SA.."
create_service_account

echo "Creating service account secret.."
create_service_account_secret

echo "Creating Cluster Issuer.."
create_cluster_issuer

echo "Creating Certificate.."
create_certificate

echo "Update Knative Ingress Gateway.."
update_knative_ingress_gateway
