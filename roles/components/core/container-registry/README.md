# Container Registry

Deploys the Harbor container registry.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#container-registry).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 6 | 8Gi | *75Gi |

*= Large storage requirement due to TAC dependency


## Variables


### Default Variables

The following variables are defaulted at `common/vars/main.yaml` and are included via `.dependencies.yaml`. They
do not need to be explicitly specified.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_container_registry.namespace | Namespace for Harbor components | tanzu-container-registry | string | yes |
| tanzu_container_registry.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-container-registry" | string | yes |
| tanzu_container_registry.admin_username | Harbor admin user | admin | string | yes |
| tanzu_container_registry.admin_password | Harbor admin password | tanzu | string | yes |
| tanzu_container_registry.{postgres,redis,registry,jobservice,chartmuseum,trivy}.size | Size of the persitent volume to use for the component | Varies | string | yes |
| tanzu_container_registry.{postgres,redis,core,notary.server,notary.signer,clair,portal,registry,jobservice,chartmuseum,trivy}.size | Resource name to use for components | harbor-{component_name} | string | yes |
| tanzu_container_registry.{postgres,redis,core,notary.server,notary.signer,clair,portal,registry,jobservice,chartmuseum,trivy}.resources | Default resource allocation for components | Varies.  See `common/defaults/main.yaml` | dict | yes |
| tanzu_container_registry.{postgres,redis,core,notary.server,notary.signer,clair,portal,registry,jobservice,chartmuseum,trivy}.replicas | Number of component pods to run | 1 | integer | yes |
| tanzu_container_registry.{postgres,redis,core,notary.server,notary.signer,clair,portal,registry,jobservice,chartmuseum,trivy}.replicas | Image to use for components | Varies | string | yes |
| tanzu_container_registry.core.dns | DNS for Harbor | registry.{{ tanzu_ingress_domain }} | string | yes |
| tanzu_container_registry.notary.dns | DNS for Notary | notary.{{ tanzu_ingress_domain }} | string | yes |
| tanzu_container_registry.{core,notary}.{ca_name,tls_name,ingress_tls_name,tls_provider} | Certificates for components | Varies | string | yes |
| tanzu_container_registry.demo_project.name | Name of the demo project to create | tanzu-demo | string | yes |
| tanzu_container_registry.demo_project.public | Whether the demo project is public or not | true | boolean | yes |
| tanzu_container_registry.demo_robot_account.name | Name of the demo robot account to create | deployer | string | yes |


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

* storage
* security
* ingress


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=container-registry make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=container-registry make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=container-registry make demo.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=container-registry make demo.test.role
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=container-registry make clean.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=container-registry make clean.test.role
```


## Accessing the Registry

After the role has been deployed, you can access your registry web UI with the following information:

* **URL:** `https://{{ tanzu_container_registry.core.dns }}.{{ tanzu_ingress_domain }}` (e.g. registry.mydomain.com)
* **Username:** `{{ tanzu_container_registry.admin_username }}` (default: admin)
* **Password:** `{{ tanzu_container_registry.admin_password }}` (default: tanzu)

You can also login with a Docker client and push/pull images:

```bash
docker login registry.mydomain.com
docker tag myproject/myimage:mytag registry.mydomain.com/myproject/myimage:mytag
docker push registry.mydomain.com/myproject/myimage:mytag
```


## Author(s)

[Andriy Kalashnykov](mailto:akalashnykov@vmware.com)
[Dustin Scott](mailto:sdustin@vmware.com)
