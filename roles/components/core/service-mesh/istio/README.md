# Service Mesh

An Ansible role to deploy an istio operator and istio control plane to a kubernetes cluster.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 4 | 4Gi | N/A |


## Variables


### Default Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_mesh.namespace | The kubernetes namespace to deploy the istio operator and components to. | "tanzu-mesh" | string | yes |
| tanzu_mesh.staging_dir | Directory to write templated kubernetes manifests to. | "{{ rpk_staging_dir }}/tanzu-mesh" | string | yes |
| tanzu_mesh.image_repo | Repository where the Istio images are stored | "projects.registry.vmware.com/rpk" | string | yes |
| tanzu_mesh.image_repo | Tag for all Istio images | "1.7.0" | string | yes |
| tanzu_mesh.policy_profile | The profile to use for the Istio deployment | "demo" | string | yes |
| tanzu_mesh.pod_labels | The labels to apply to all Istio pods | See `common/vars/main.yaml` | dict | yes |
| tanzu_mesh.operator.name | Name to create the istioOperator CRD object as. | "istiocontrolplane" | string | yes |
| tanzu_mesh.operator.image | Name of the image to pull from `tanzu_mesh.image_repo`. | "istio-operator" | string | yes |
| tanzu_mesh.operator.{service_account,clusterrole,clusterrolebinding} | RBAC names for Istio objects. | See `common/vars/main.yaml` | string | yes |
| tanzu_mesh.operator.resources | Resources to apply to Istio operator | See `common/vars/main.yaml` | dict | yes |
| tanzu_mesh.{operator,istiod,ingress_gateway,egress_gateway}.pod_labels | The labels to apply to Istio components | See `common/vars/main.yaml` | dict | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=service-mesh/istio make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=service-mesh/istio make deploy.test.role
```


## Demonstrating

**NOTE:** no demo role is avaiable for Istio at this time.

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=service-mesh/istio make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=service-mesh/istio make clean.role
```


## Author(s)
[Andrew J. Huffman](mailto:ahuffman@vmware.com)
