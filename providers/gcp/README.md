## Google Cloud Platform Provider

Hacky & bare implementation with shell scripts.

/!\ The service account associated with the credentials json must have the "container.clusterRoleBindings.create" permission.

You can add it by running:

```
gcloud projects add-iam-policy-binding <project> \
  --member=serviceAccount:<service-account>@developer.gserviceaccount.com \
  --role=roles/container.admin
```

# Status

| Command | Status |
| - | - |
| setup | ✔ |
| delete | ✔ |
| get | ✔ |
| app deploy | ✔ |
| app delete | ✔ |
| dev start | ✔ |
| dev stop | ✔ |

# Configuration

You can override the following settings by adding them as flags of the `kuda setup` command (ex: `kuda setup gcp ... --gcp_cluster_name=mycluster`).

| Parameter | Default | Description |
| - | - | - |
| `gcp_project_id` | None (Required) | The GCP Project ID |
| `gcp_credentials` | None (Required) | Path to the GCP Credential JSON file |
| `gcp_cluster_name` | kuda | The new or existing cluster name |
| `gcp_compute_zone` | us-central1-a | The GCP compute zone |
| `gcp_machine_type` | n1-standard-4 | Default machine type for the nodes (Only evaluated during `setup`)|
| `gcp_pool_num_nodes` | 1 | Default number of nodes of the GPU pool (Only evaluated during `setup`) |
| `gcp_gpu` | k80 | The default GPU to use. (Only evaluated during `setup`) |
| `gcp_use_preemptible` | false | Wether or not the GPU nodes should be preemptible. (Only evaluated during `setup`) |


# Limitations

- The "Compute Engine API - GPUs (all regions)" quota must be requested manually [here](<https://console.cloud.google.com/iam-admin/quotas?metric=GPUs%20(all%20regions)>)
- By default, the system node and the load balancer will be kept on, incuring charges of about 45€ per months. You can manually scale down the system node to 0 to temporarily stop its associated charges or run `kuda delete` to completely delete the cluster.
- Currently the load balancer doesn't get deleted when you delete the cluster. Make sure to delete it manually [here](https://console.cloud.google.com/net-services/loadbalancing/loadBalancers/list) after deleting a cluster to avoid extra costs.
- You can find a list of parameters that you can override in `.config.sh`
