#!/bin/bash

set -e

source $KUDA_CMD_DIR/.config.sh

# Get cluster's credentials to use kubectl.
gcloud container clusters get-credentials $KUDA_GCP_CLUSTER_NAME

# TODO: Increase the number of GPU nodes by 1 to speed up initialization.
# gcloud container clusters resize $KUDA_GCP_CLUSTER_NAME \
#   --node-pool $KUDA_DEFAULT_GPU \
#   --num-nodes 1 \
#   --quiet

# TODO: Support multiple sessions.
# tmp_uuid=$(od -x /dev/urandom | head -1 | awk '{print $2$3}')
# KUDA_DEV_APP_NAME=kuda-dev-$tmp_uuid

# exit 1

# Launch.
# TODO: materialize as external yaml file and customize
# with Kustomize
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: $KUDA_DEV_APP_NAME
  labels:
    app: $KUDA_DEV_APP_NAME
spec:
  ports:
  - name: http
    port: 8000
    targetPort: 80
  selector:
    app: $KUDA_DEV_APP_NAME
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $KUDA_DEV_APP_NAME
spec:
  replicas: 1
  selector:
    matchLabels:
      app: $KUDA_DEV_APP_NAME
  template:
    metadata:
      labels:
        app: $KUDA_DEV_APP_NAME
    spec:
      containers:
        - name: $KUDA_DEV_APP_IMAGE_NAME
          image: $KUDA_DEV_APP_IMAGE
          resources:
            limits:
              nvidia.com/gpu: 1
          ports:
          - containerPort: 80
          env:
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: KUDA_PROVIDER
              value: gcp
            - name: KUDA_GCP_PROJECT_ID
              value: $KUDA_GCP_PROJECT_ID
            - name: KUDA_GCP_CREDENTIALS
              value: $KUDA_GCP_CREDENTIALS
          volumeMounts:
            - name: secret
              readOnly: true
              mountPath: "/secret"
      volumes:
        - name: secret
          secret:
            secretName: $(basename $KUDA_GCP_CREDENTIALS)
EOF

# Setup Istio Ingress gateway.
cat <<EOF | kubectl apply -f -
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: $KUDA_DEV_APP_NAME-gateway
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "$KUDA_DEV_APP_NAME.example.com"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: $KUDA_DEV_APP_NAME
spec:
  hosts:
  - "$KUDA_DEV_APP_NAME.example.com"
  gateways:
  - $KUDA_DEV_APP_NAME-gateway
  http:
  - match:
    - uri:
        prefix: /
    route:
    - destination:
        port:
          number: 8000
        host: $KUDA_DEV_APP_NAME
EOF

# Initialize ksync.
echo "Setting up Ksync.."
ksync init

# Patch Ksync DaemonSet to be scheduled on gpu nodes.
echo "Patching Ksync DaemonSet.."
kubectl -n kube-system patch daemonset ksync --type merge -p '
{
   "spec": {
      "template": {
         "spec": {
            "tolerations": [
               {
                  "effect": "NoSchedule",
                  "key": "nvidia.com/gpu",
                  "operator": "Exists"
               }
            ]
         }
      }
   }
}
'

# Wait for readiness.
# TODO: Show special message to indicate installation of Nvidia Driver.
echo "Starting...  "
i=1
sp="⣷⣯⣟⡿⢿⣻⣽⣾"
while true; do
  printf "\b${sp:i++%${#sp}:1}"
  min=1
  status=$(kubectl get deployment $KUDA_DEV_APP_NAME -o jsonpath={.status.availableReplicas})
  if [ ! -z "$status" ]; then
    if [ "$status" -ge 1 ]; then
      break
    fi
  fi
  sleep 0.1
done
echo

# Startup local client.
ksync watch --daemon=true

# Watch file changes.
ksync create \
  --selector=app=$KUDA_DEV_APP_NAME \
  --reload=false \
  $KUDA_DEV_SYNC_PATH \
  $KUDA_DEV_SYNC_PATH

# Allow some time for ksync to synchronize the base folder.
sleep 3

# Print IP address
ingress_host=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
ingress_port=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].port}')
echo
echo "$ingress_host:$ingress_port"
echo

# Launch remote shell.
pod_name=$(kubectl get pods -o name | grep -m1 $KUDA_DEV_APP_NAME | cut -d'/' -f 2)
kubectl exec -it $pod_name -- /bin/bash
