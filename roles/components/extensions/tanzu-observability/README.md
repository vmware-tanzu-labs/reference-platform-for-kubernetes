# Tanzu Observability

The following role will register your TKG cluster with Tanzu Observability.  For official docs,
please visit https://docs.wavefront.com/kubernetes.html.

The registration and integration of Tanzu Observability can be performed in one of two ways, either via the Kubernetes manifests detailed below or via native Tanzu Mission Control (TMC) enabled integrations once an RPK cluster has been successfully attached to or provisioned via TMC.
Further details on TMC integrations can be found in https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-6DFB45C0-A741-4972-95EA-59B7AE581FE8.html

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Origination of Manifests

The templates were created with the following commands from https://docs.wavefront.com/kubernetes.html :

```bash
helm repo add wavefront https://wavefronthq.github.io/helm/
helm repo update
helm pull wavefront/wavefront
helm template wavefront ./ > ../../app-wavefront.yaml.j2
```

## TMC API Usage

The TMC API scheme and calls used in this role are further documented in https://code.vmware.com/apis/1079/tanzu-mission-control

## Resoure Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 2000m | 2Gi | N/A - no persistent storage required |

**NOTE:** above sizing requirements will adjust based on cluster size due to TO
components utilizing a DaemonSet.

## Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_observability.namespace | Namespace for extensions/tanzu-observability components | "extensions/tanzu-observability" | string | yes |
| tanzu_observability.staging_dir | Local directory to write the staging manfiests to | "/tmp/staging/tanzu-observability" | string | yes |
| tanzu_observability.tmc_access_token_url | URL to pull the TMC access token from | "https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize" | string | yes |
| tanzu_observability.tmc_org_url | Base URL of the TMC organization to integrate with | "https://{{ tanzu_mission_control_org_name }}.tmc.cloud.vmware.com" | string | yes |
| tanzu_observability.tmc_api_version | API Version of TMC to use | "v1alpha" | string | yes |
| tanzu_observability.use_tmc | Use Tanzu Mission Control to integrate cluster with Tanzu Observability | false | boolean | yes |
| tanzu_observability.collector.image | The image to use for Tanzu Observability collector | "wavefronthq/wavefront-kubernetes-collector" | string | yes |
| tanzu_observability.collector.image_tag | The image tag to use for Tanzu Observability collector | "1.2.3" | string | yes |
| tanzu_observability.proxy.image | The image to use for Tanzu Observability proxy | "wavefronthq/proxy" | string | yes |
| tanzu_observability.proxy.image_tag | The image tag to use for Tanzu Observability proxy | "9.2" | string | yes |
| tanzu_observability.proxy.replicas | The number of replicas to use for Tanzu Observability proxy | 1 | integer | yes |
| tanzu_observability.workload_tenancy.enabled | Whether to use the `workload-tenancy` module to provide custom namespaces | false | boolean | yes |


## Additional Variables

The following variables must be set for proper operatiion of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_observability_api_token | API Token to use for Tanzu Observability | "" | string | yes |
| tanzu_observability_org_name | Organization name (from TO URL - e.g vmware.wavefront.com = vmware) | "" | string | yes |
| tanzu_mission_control_org_name | Name of the org (from TMC URL) | - | string | no |
| tanzu_mission_control_api_token | Token which is authorized for TMC from cloud.vmware.com console | - | string | no |
| tanzu_cluster_name | The name of the cluster which will be attached to Tanzu Observability | "" | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role:

```bash
# deploy the role
ROLE=extensions/tanzu-observability make deploy.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

Demonstrating this role simply prints out the access information for Tanzu Observability.

In order to demonstrate the role:

```bash
ROLE=extensions/tanzu-observability make demo.role
```

Sample output:

```bash
...
running rpk deployment tool...

PLAY [all] ***********************************************************************************************************

TASK [extensions/tanzu-observability/demo : print the dashboard access information] **********************************
ok: [aws_tanzu_cluster1] =>
  msg: You can access the Tanzu Observability Dashboard at URL http://vmware.wavefront.com
```


## Cleaning

From the root of the repo:

```bash
ROLE=extensions/tanzu-observability make clean.role
```

## Watch a Demo

[Tanzu Observability Demo](https://drive.google.com/file/d/1HuGRuOpSWmdMmGKhkru8nTpiz0a8Q7bO/view?usp=sharing)

## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)

[Paul Wiggett](mailto:pwiggett@vmware.com)
