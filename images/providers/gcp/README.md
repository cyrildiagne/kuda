## Google Cloud Platform Provider

Hacky & bare implementation with shell scripts.

# Limitations

- The "Compute Engine API - GPUs (all regions)" quota must be requested manually [here](https://console.cloud.google.com/iam-admin/quotas?metric=GPUs%20(all%20regions))
- Currently the load balancer doesn't get deleted when you delete the cluster. Make sure to delete it manually [here](https://console.cloud.google.com/net-services/loadbalancing/loadBalancers/list) after deleting a cluster to avoid extra costs.
- You can find a list of parameters that you can override in `.config.sh`