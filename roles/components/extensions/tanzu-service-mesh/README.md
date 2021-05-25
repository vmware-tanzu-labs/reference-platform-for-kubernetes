# Tanzu Service Mesh

The following role will register your TKG cluster with Tanzu Service Mesh.  For official documentation,
please visit https://docs.vmware.com/en/VMware-Tanzu-Service-Mesh/index.html.

The registration and integration of Tanzu Service Mesh can be performed in one of two ways, either via the Kubernetes manifests detailed below or via native Tanzu Mission Control (TMC) enabled integrations once an RPK cluster has been successfully attached to or provisioned via TMC.
Further details on TMC integrations can be found in https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-6DFB45C0-A741-4972-95EA-59B7AE581FE8.html

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)

## Origination of Manifests

The templates were created with the following commands:

```bash
curl https://prod-2.nsxservicemesh.vmware.com/cluster-registration/k8s/operator-deployment.yaml -o operator-deployment.yaml.j2
```

## TMC API Usage

The TMC API scheme and calls used in this role are further documented in https://code.vmware.com/apis/1079/tanzu-mission-control

## Resoure Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 25 | 30Gi | N/A - no persistent storage required |

## Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_service_mesh.namespaces.tsm | Namespace for extensions/tanzu-service-mesh operator components | "vmware-system-tsm" | string | yes |
| tanzu_service_mesh.namespaces.istio | Namespace for extensions/tanzu-service-mesh istio components | "istio-system" | string | yes |
| tanzu_service_mesh.staging_dir | Local directory to write the staging manfiests to | "/tmp/staging/tanzu-service-mesh" | string | yes |
| tanzu_service_mesh.tmc_access_token_url | URL to pull the TMC access token from | "https://console.cloud.vmware.com/csp/gateway/am/api/auth/api-tokens/authorize" | string | yes |
| tanzu_service_mesh.tmc_org_url | Base URL of the TMC organization to integrate with | "https://{{ tanzu_mission_control_org_name }}.tmc.cloud.vmware.com" | string | yes |
| tanzu_service_mesh.tmc_api_version | API Version of TMC to use | "v1alpha" | string | yes |
| tanzu_service_mesh.use_tmc | Use Tanzu Mission Control to integrate cluster with Tanzu Service Mesh | false | boolean | yes |
| tanzu_service_mesh.images.tsm_agent_operator.image | The image to use for Tanzu Service Mesh Operator | "284299419820.dkr.ecr.us-west-2.amazonaws.com/tsm-agent-operator" | string | yes |
| tanzu_service_mesh.images.tsm_agent_operator.image_tag | The image tag to use for Tanzu Service Mesh Operator | "v1.2.7" | string | yes |
| tanzu_service_mesh.images.operator_ecr.image | The image to use for Tanzu Service Mesh operator-ecr-read-only--renew-token | "docker.io/vmwareallspark/photon-kubectl" | string | yes |
| tanzu_service_mesh.images.operator_ecr.image_tag | The image tag to use for Tanzu Service Mesh operator-ecr-read-only--renew-token | "1.15.0" | string | yes |
| tanzu_service_mesh.workload_tenancy.enabled | Whether to use the `workload-tenancy` module to provide custom namespaces | false | boolean | yes |

## Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_service_mesh_api_token | Cloud Services Portal Token to use for Tanzu Service Mesh cluster onboarding. See [Tanzu Service Mesh - Authentication with the Tanzu Service Mesh REST API Documentation](https://docs.vmware.com/en/VMware-Tanzu-Service-Mesh/services/api-programming-guide/GUID-FED8E849-B3C3-49ED-9FDB-1317CFFF3141.html) and follow step 1 of the Procedure section. | "" | string | yes |
| tanzu_service_mesh_url | URL to your Tanzu Service Mesh instance (e.g. https://prod-2.nsxservicemesh.vmware.com) | "" | string | yes |
| tanzu_mission_control_org_name | Name of the org (from TMC URL) | - | string | no |
| tanzu_mission_control_api_token | Token which is authorized for TMC from cloud.vmware.com console | - | string | no |
| tanzu_cluster_name | Name of your Kubernetes cluster | "" | string | yes |

## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE

## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role:

```bash
# deploy the role
ROLE=extensions/tanzu-service-mesh make deploy.role
```

## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=extensions/tanzu-service-mesh make demo.role
```

Sample output:

```bash
...
running RPK deployment tool...

PLAY [all] ***********************************************************************************************************

TASK [extensions/tanzu-service-mesh/demo : print the dashboard access information] **********************************
ok: [aws_tanzu_cluster1] =>
  msg: You can access the Tanzu Service Mesh Dashboard at URL https://console.cloud.vmware.com
```

## Cleaning

From the root of the repo:

```bash
ROLE=extensions/tanzu-service-mesh make clean.role
```

## Author(s)
[Andrew J. Huffman](mailto:ahuffman@vmware.com)

[Paul Wiggett](mailto:pwiggett@vmware.com)
