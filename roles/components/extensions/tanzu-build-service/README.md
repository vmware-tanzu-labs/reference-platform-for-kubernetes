# Tanzu Build Service

This role deploys Tanzu Build Service.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 4 | 12Gi | 5 to 10 GB |


## Variables

### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name                                     | Description                                                               | Default Value                      | Variable Type | Required |
|---------------------------------------------------|---------------------------------------------------------------------------|------------------------------------|---------------|----------|
| tanzu_build_service.namespace                     | Namespace for build-service components                                    | build-service ( do not change )    | string        | yes      |
| tanzu_build_service.staging_dir                   | Local directory to write the staging manfiests to                         | "/tmp/staging/tanzu-build-service" | string        | yes      |
| tanzu_build_service.namespace_kpack               | Namespace for kpack components                                            | kpack ( do not change )            | string        | yes      |
| tanzu_build_service.registry.project.project_name | Harbor Registry name to use                                               | build-service                      | string        | yes      |
| tanzu_build_service.registry.source_url           | Source Tanzu Build Service Registry URL ( do not change )                 | registry.pivotal.io                | string        | yes      |
| tanzu_build_service.workload_tenancy.enabled      | Whether to use the `workload-tenancy` module to provide custom namespaces | false                              | boolean       | yes      |

### Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name                | Description                                                | Default Value                      | Variable Type | Required |
| ---------------------------- | ---------------------------------------------------------- | ---------------------------------- | ------------- | -------- |
| tanzu_kubectl_context        | Name of context to use for connection to Kubernetes        | -                                  | string        | yes      |
| tanzu_build_service_username | Username for registry.pivotal.io                           | -                                  | string        | yes      |
| tanzu_build_service_password | Password for registry.pivotal.io                           | -                                  | string        | yes      |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* storage
* security
* ingress
* container-registry


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=extensions/tanzu-build-service make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=extensions/tanzu-build-service make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

Brief description of what the demonstration does here.

In order to demonstrate the role:

```bash
ROLE=extensions/tanzu-build-service make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=extensions/tanzu-build-service make clean.role
```

## Watch a Demo

[Tanzu Build Service Demo](https://drive.google.com/file/d/1il9bai8dXSwCxFg0sMaLgVlZF6rjbcZt/view?usp=sharing)

## Author(s)
[Robin Foe](mailto:rfoe@vmware.com)
