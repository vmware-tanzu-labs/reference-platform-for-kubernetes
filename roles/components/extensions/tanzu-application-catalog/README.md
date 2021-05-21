# Tanzu Application Catalog

Creates a registry endpoint and a replication job in the Harbor registry to sync images from upstream Tanzu Application Catalog.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 2 | 4Gi | 50GB |

**NOTE:** storage sizing is dependent upon `{{ tanzu_container_registry.registry.size }}`


## Variables


### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_application_catalog.namespace | Namespace for extensions/tanzu-application-catalog components | "tanzu-application-catalog" | string | yes |
| tanzu_application_catalog.staging_dir | Local directory to write the staging manfiests to | "/tmp/staging/extensions/tanzu-application-catalog" | string | yes |
| tanzu_application_catalog.registry_endpoint | Destination registry for images (defaults to local Harbor instance) | "{{ tanzu_container_registry.core.dns }}" | string | no |
| tanzu_application_catalog.remote_registry_name | Remote registry name (to sync from) | "tac_registry" | string | yes |
| tanzu_application_catalog.replication_policy_name | Name for the replication policy | "tac" | string | yes |
| tanzu_application_catalog.cron_schedule | Sync execution schedule in cron format | "0 40 4 * * 1" | string | yes |
| tanzu_application_catalog.app_list | List of images to sync from upstream repository (defaults to dynamically discovering the app list) | [] | list | yes |
| tanzu_application_catalog.override_map | List of overrides from app_list discovery to override when `bitnami_app_id` doesn't match the image name | See `common/vars/main.yaml` | dict | yes |


### Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_application_catalog_username  | Username of the robot account with access the remote registry (e.g. `robot$my-account`) | "" | string | yes |
| tanzu_application_catalog_password  | Password to access the remote registry | "" | string | yes |
| tanzu_application_catalog_api_token  | API Token to list TAC artifacts | "" | string | yes |


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
ROLE=extensions/tanzu-application-catalog make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=extensions/tanzu-application-catalog make deploy.test.role
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=extensions/tanzu-application-catalog make clean.role
```


## Author(s)
[Farid Saad](mailto:fsaad@vmware.com)
