# Secret Management / Hashicorp Vault

An Ansible role to deploy Hashicorp Vault as an external kubernetes secret management tool.

The role currently assumes that the inventory target will be a node with kubectl configured with administrator rights.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md#vault-integration).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 4 | 4Gi | 10Gi |


## Variables


### Default Variables

Default variables are located at `common/vars` within this role:

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_secrets.namespace | Namespace to direct `kubectl` commands and separate, but dependent on how the Kubernetes manifests have been configured. | "tanzu-secrets" | string | yes |
| tanzu_secrets.dns | DNS name to use for Vault | secrets | string | yes |
| tanzu_secrets.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/{{ tanzu_secrets.namespace }}" | string | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.image | Image to use for Vault/Injector | Varies | string | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.image_tag | Image Tag to use for Vault/Injector | Varies | string | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.replicas | Number of replicas to use for Vault/Injector | 1 | integer | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.resource_name | Name of service account for Vault/Injector | Varies | string | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.resource_name | Name of cluster role for Vault/Injector | Varies | string | yes |
| tanzu_secrets.hashicorp_vault}.resource_name | Name of cluster role binding for Vault/Injector | Varies | string | yes |
| tanzu_secrets.hashicorp_vault.secret_name | Name of secret to store root token and unseal keys | hashicorp-vault-secrets | string | yes |
| tanzu_secrets.hashicorp_vault.engine_backend | Name of the Vault engine backend to use | kv-v2 | string | yes |
| tanzu_secrets.hashicorp_vault.engine_name | Name of the vault engine to use | tanzu | string | yes |
| tanzu_secrets.hashicorp_vault.policy_name | Name of the vault policy to use | tanzu | string | yes |
| tanzu_secrets.hashicorp_vault.role_name | Name of the vault role to use  | tanzu | string | yes |
| tanzu_secrets.{hashicorp_vault_injector,hashicorp_vault}.resources | Default resource allocation vault components | Varies.  See `common/vars/main.yaml` | dict | yes |
| tanzu_secrets.hashicorp_vault.demo.app_name | Name of the demo app to deploy with `make demo.role` | tanzu-secrets-demo-app | string | yes |


### Additional Variables

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
ROLE=secret-management/hashicorp-vault make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=secret-management/hashicorp-vault make deploy.test.role
```


## Verifying

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=secret-management/hashicorp-vault make demo.role
```

Sample output:

```bash
...
ok: [vmware_tanzu_cluster1] => {
    "msg": [
        "You can access the Vault UI at URL http://secrets.vsphere.ihazcloudthings.xyz",
        "Root Token: s.dgAE70V2SHoH5u0ypEidyKVL",
        "Unseal Keys: ab9df98486960260b38ef28d28f26f269bf4d989bd992b501dbd8016dd6e3384ad 533a2c992019343bdd5008ce614a793857840af8435a2c418cc0a379f1de8eb547 7904f41a26453b2ce735fcfed4468232679561626d11167d5308c3f843b7774066 de2b7d5f4001760259751f98321295009df13b57361f02f4baa720686897e4fe35 bf26f480b4c518d0f7412c8c624d77ab05acb8fbb1048b10a0026bebb447e444c2",
        "View the mounted Vault secrets via demo app pod at /vault: tanzu-secrets-demo-app"
    ]
}
...
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=secret-management/hashicorp-vault make clean.role
```


## Author(s)
[Andrew J. Huffman](mailto:ahuffman@vmware.com)
