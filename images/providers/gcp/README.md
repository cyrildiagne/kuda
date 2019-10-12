## Google Cloud Platform Provider

Hacky & bare implementation with shell scripts.

/!\ The service account associated with the credentials json must have the "container.clusterRoleBindings.create" permission.

You can add it by running:

```
gcloud projects add-iam-policy-binding <project> \
  --member=serviceAccount:<service-account>@developer.gserviceaccount.com \
  --role=roles/container.admin
```

# Limitations

- The "Compute Engine API - GPUs (all regions)" quota must be requested manually [here](<https://console.cloud.google.com/iam-admin/quotas?metric=GPUs%20(all%20regions)>)
- By default, the system node and the load balancer will be kept on, incuring charges of about 45€ per months. You can manually scale down the system node to 0 to temporarily stop its associated charges or run `kuda delete` to completely delete the cluster.
- Currently the load balancer doesn't get deleted when you delete the cluster. Make sure to delete it manually [here](https://console.cloud.google.com/net-services/loadbalancing/loadBalancers/list) after deleting a cluster to avoid extra costs.
- You can find a list of parameters that you can override in `.config.sh`

# Development

