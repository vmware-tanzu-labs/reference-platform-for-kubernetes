# Workload Tenancy

Worload Tenancy implements best practices as defined from Tanzu Developer Center (TDC).  Please visit https://tanzu.vmware.com/developer/guides/kubernetes/workload-tenancy/
for more information.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#workload-tenancy).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 50m | 64Mi | N/A - no persistent storage required |

## Default Variables

The following variables are defaulted at `common/vars/main.yaml` and are included via `.dependencies.yaml`. They
do not need to be explicitly specified.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_workload_tenancy.namespace | Defines the name for the workload-tenancy namespace | "tanzu-workload-tenancy" | string | yes |
| tanzu_workload_tenancy.demo_namespace | Defines the name for the workload-tenancy namespace used by `make demo.role` | "tanzu-demo" | string | yes |
| tanzu_workload_tenancy.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-workload-tenancy" | string | yes |
| tanzu_workload_tenancy.namespace_operator.service_account | Service account used to run the namespace-operator | "namespace-operator" | string | yes |
| tanzu_workload_tenancy.namespace_operator.clusterrole | Cluster Role used by service account which runs the namespace-operator | "namespace-operator-clusterrole" | string | yes |
| tanzu_workload_tenancy.namespace_operator.clusterrolebinding | Cluster Role Binding used by service account/role which runs the namespace-operator | "namespace-operator-clusterrolebinding" | string | yes |
| tanzu_workload_tenancy.namespace_operator.image | namespace-operator image | "scottd018/namespace-operator" | string | yes |
| tanzu_workload_tenancy.namespace_operator.image_tag | namespace-operator image tag | "v0.0.1beta" | string | yes |
| tanzu_workload_tenancy.namespace_operator.replicas | namespace-operator replica count | 2 | integer | yes |
| tanzu_workload_tenancy.namespace_operator.resources | Normal Kubernetes resource construct defining resource requirements | See `common/vars/main.yaml` | dict | yes |

## Additional Variables

The following variables are required for proper role operation and should be specified as part
of the global inventory, or other mechanism.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=workload-tenancy make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=workload-tenancy make deploy.test.role
```

## Validating

For a module that has validation tests, run them with:

    $ ROLE=workload-tenancy make validate.role

To view the status of the validation test:

    $ make validate.status

Once status is `complete`, If test result is `passed`, clean up the sonobuoy
assets as follows:

    $ make validate.clean

See docs/VALIDATION.md for further info regarding validation tests.

## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.  This demo simply
does the following:

- Creates a TanzuNamespace resource, which in turn creates:
  - Namespace
  - LimitRange in the Namespace
  - ResourceQuota in the Namespace
- Creates a pod which utilizes defaults by the LimitRange created
- Creates a pod which fails to launch by exceeding LimitRange CPU limits
- Creates a pod which fails to launch by exceeding LimitRange Memory limits

```bash
ROLE=workload-tenancy make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=workload-tenancy make clean.role
```


## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
[Rich Lander](mailto:landerr@vmware.com)
