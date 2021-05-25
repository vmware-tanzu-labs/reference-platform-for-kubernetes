# Tanzu Mission Control

This role will do the following to your TKG Cluster:

- Create a Cluster Group in TMC (if non-existing)
- Create a Cluster in a Cluster Group
- Attach a Cluster

See https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/index.html for more information on
Tanzu Mission Control.


## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Resoure Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

**NOTE:** this is only an estimate.  Please refer to official product documentation for official details (https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/index.html).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 500m | 1Gi | N/A - no persistent storage required |

## Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_mission_control.namespace | Namespace for TMC components | "vmware-system-tmc" | string | yes |
| tanzu_mission_control.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-mission-control" | string | yes |
| tanzu_mission_control.access_token_url | URL to pull the access token from | "https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize" | string | yes |
| tanzu_mission_control.org_url | Base URL of the organization to integrate with | "https://{{ tanzu_mission_control_org_name }}.tmc.cloud.vmware.com" | string | yes |
| tanzu_mission_control.api_version | API Version of TMC to use | "v1alpha" | string | yes |

## Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_mission_control_org_name | Name of the org (from TMC URL) | - | string | yes |
| tanzu_mission_control_api_token | Token which is authorized for TMC from cloud.vmware.com console | - | string | yes |
| tanzu_mission_control_cluster_group | Name of the cluster group to register into (created if non-existent) | - | string | yes |
| tanzu_cluster_name | DNS friendly name (no underscores or special characters) of the cluster to register into TMC | - | string | yes |


## tanzu-mission-control Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role:

```bash
# only run this if you need to build a new image or haven't built one yet
make build

# deploy the role
ROLE=extensions/tanzu-mission-control make deploy.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

**NOTE:** no demo for TMC at this time.


## Cleaning

The following will "detach" the cluster from Tanzu Mission Control and remove all associated objects:

From the root of the repo:

```bash
ROLE=product/tanzu-mission-control make clean.role
```

## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
