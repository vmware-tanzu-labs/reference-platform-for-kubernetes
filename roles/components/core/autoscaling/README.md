# Autoscaling

Autoscaling is the home for any autoscaling functionality used on the RPK platform. Currently, it includes the [Vertical Pod Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler) (VPA).  In order to use, create a `VerticalPodAutoscaler` resource to configure vertical autoscaling behavior for your workload.  An example manifest can be found [here](https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler#example-vpa-configuration).

## Default Variables

The following variables are defaulted at `common/vars/main.yaml` and are included via `meta/main.yaml`. They
do not need to be explicitly specified.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_autoscaling.namespace | Defines the name for the autoscaling namespace | "tanzu-autoscaling" | string | yes |
| tanzu_autoscaling.demo_namespace | Defines the name for the autoscaling namespace used by `make demo.role` | "tanzu-demo" | string | yes |
| tanzu_autoscaling.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-autoscaling" | string | yes |
| tanzu_workload_tenancy.vpa.version | Version of the Vertical Pod Autoscaler to use | "0.9.0" | string | yes |

## Additional Variables

The following variables are required for proper role operation and should be specified as part
of the global inventory, or other mechanism.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |

## Dependencies

Also see `meta/main.yaml` to view role dependencies which are run when running the role
independently.

* security
* monitoring

## Deploying

**NOTE:** roles from `meta/main.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=autoscaling make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=autoscaling make deploy.test.role
```

## Validating

For a module that has validation tests, run them with:

    $ ROLE=autoscaling make validate.role

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
- Creates a sample nginx Deployment and Service
- Creates a VerticalPodAutoscaler in `Initial` mode for the nginx sample workload

```bash
ROLE=autoscaling make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=autoscaling make clean.role
```

## Author(s)
[Rich Lander](mailto:landerr@vmware.com)
