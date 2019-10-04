#!/bin/bash

source $KUDA_CMD_DIR/.config.sh

sleep 1

while true; do
  echo
  read -p "DANGER! This will delete the cluster $KUDA_GCP_CLUSTER_NAME. Are you sure (y/n)? " yn
  case $yn in
  [Yy]*) break ;;
  [Nn]*)
    echo "Cancelling..."
    exit
    ;;
  *) echo "y/n" ;;
  esac
done

echo "Deleting cluster... (this may take a while)"

# Delete cluster.
gcloud container clusters delete $KUDA_GCP_CLUSTER_NAME --quiet

# Delete the load balancer.
while true; do
  echo
  read -p "DANGER! Do you want to delete the all the backend-services of project $KUDA_GCP_PROJECT? (y/n)" yn
  case $yn in
  [Yy]*) break ;;
  [Nn]*)
    echo "Leaving the backend-services. Check manually for orphaned resources."
    exit
    ;;
  *) echo "y/n" ;;
  esac
done

echo "Deleting the load balancer & backend services..."

# Delete fowarding rules.
gcloud compute forwarding-rules delete https-content-rule --quiet
# Delete the global external IP addresses.
gcloud compute addresses delete lb-ipv4-1 --quiet
gcloud compute addresses delete lb-ipv6-1 --quiet
# Delete the target proxy.
gcloud compute target-https-proxies delete https-lb-proxy --quiet
# Delete the SSL certificate.
gcloud compute ssl-certificates delete www-ssl-cert --quiet
# Delete the URL map.
gcloud compute url-maps delete web-map --quiet
# Delete the backend services.
gcloud compute backend-services delete web-backend-service --quiet
gcloud compute backend-services delete video-backend-service --quiet
# Delete the health check.
gcloud compute health-checks delete http-basic-check --quiet
