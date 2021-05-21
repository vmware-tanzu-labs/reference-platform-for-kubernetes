# Dashboard

Deploy the Octant dashboard.

**NOTE:** at this time deploying the Octant dashboard within the Kubernetes cluster itself is unsupported.  This
should be used for demo purposes only.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 1 | 2Gi | N/A - no persistent storage required |


## Variables

### Default Variables

The following variables are defaulted at `common/vars/main.yaml` and are included via `.dependencies.yaml`. They
do not need to be explicitly specified.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_dashboard.namespace | Namespace for dashboard components | tanzu-dashboard | string | yes |
| tanzu_dashboard.dns | DNS name to use for Octant | dashboard.{{ tanzu_ingress_domain }} | string | yes |
| tanzu_dashboard.service_account | Service to create and use for Octant | tanzu-dashboard | string | yes |
| tanzu_dashboard.clusterrole | Cluster Role to create and use for Octant | tanzu-dashboard-clusterrole | string | yes |
| tanzu_dashboard.clusterrolebinding | Cluster Role Binding to create and use for Octant | tanzu-dashboard-clusterrolebinding | string | yes |
| tanzu_dashboard.resources | Normal Kubernetes resource construct defining resource requirements | See `common/vars/main.yaml` | dict | yes |
| tanzu_dashboard.image | Image to use for Octant | aleveille/octant-dashboard | string | yes |
| tanzu_dashboard.image_tag | Image Tag to use for Octant | latest | string | yes |
| tanzu_dashboard.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-dashboard" | string | yes |
| tanzu_dashboard.workload_tenancy.enabled | Whether to use the `workload-tenancy` module to provide custom namespaces | false | boolean | yes |


### Additional Variables

The following variables are required for proper role operation and should be specified as part
of the global inventory, or other mechanism.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_ingress_domain | Subdomain used for deployment of RPK | - | string | yes |
| ingress_ip | Ingress IP Address (from ingress role) | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* security
* ingress


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=extensions/octant make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=extensions/octant make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=extensions/octant make demo.role
```


## Cleaning

From the root of the repo:

```bash
ROLE=extensions/octant make clean.role
```


## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
