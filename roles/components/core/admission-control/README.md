# Admission Control

This role deploys admission control constructs that demonstrate both [validation](https://github.com/open-policy-agent/gatekeeper)
and [mutation of resources](https://github.com/vmware-tanzu-labs/scheduling-defaults-webhook) submitted to the Kubernetes API server.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#admission-control).

## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 1000m | 512Mi | N/A - no persistent storage required |


## Variables

### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_admission_control.namespace | Namespace for tanzu_admission_control components | "tanzu-admission" | string | yes |
| tanzu_admission_control.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-admission" | string | yes |
| tanzu_admission_control.gatekeeper.image | Image to use for gatekeeper components | "openpolicyagent/gatekeeper" | string | yes |
| tanzu_admission_control.gatekeeper.image_tag | Image tag to use for gatekeeper components | "v3.1.1" | string | yes |
| tanzu_admission_control.gatekeeper.audit.resources | Normal Kubernetes resource construct defining resource requirements | See `common/vars/main.yaml` | dict | yes |
| tanzu_admission_control.gatekeeper.controller.resources | Normal Kubernetes resource construct defining resource requirements | See `common/vars/main.yaml` | dict | yes |
| tanzu_admission_control.mutator.image | Image to use for mutator components | "joshrosso/sac" | string | yes |
| tanzu_admission_control.mutator.image_tag | Image tag to use for mutator components | "1.18" | string | yes |
| tanzu_admission_control.mutator.resources | Normal Kubernetes resource construct defining resource requirements | See `common/vars/main.yaml` | dict | yes |

### Additional Variables

The following variables must be set for proper operatiion of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* security

## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=admission-control make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=admission-control make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

**NOTE:** this is a work in progress.  See [Verification](#verification) for steps to manually demonstrate
the role.

```bash
ROLE=admission-control make demo.role
```


## Verification

### Affinity Verification

1. Deploy a workload that has no affinity set.

1. Get the deployed manifest, verify affinity has been set (by the mutating
   webhook controller.

### Root use verification

1. Deploy a workload that is not set to run as root.

1. Verify the workload is rejected.

### Required Label Verification

1. Deploy a workload without the `app` label set.

1. Verify the workload is rejected.

## Configurability

This admission stack can be easily configured, the following sections
demonstrate how custom configuration could occur.

### Adding more required labels

To add more required labels, update the constraint within `manifests/policy/labels`.

### Adding more security policy

The entire gatekeeper library is available under `manifests/policy/psp`. Apply
additional policies such as host path protection to further lock down a cluster.

### Adding or changing mutating rules

The mutating webhook source code is available at
https://github.com/vmwarepivotallabs/scheduling-defaults-webhook, the source
code can be easily modified to change the behavior.


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=admission-control make clean.role
```


## Author(s)
[Josh Rosso](mailto:jrosso@vmware.com)
