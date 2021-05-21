# Storage

Deploy storage and default classes into providers.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#storage-integration).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 2 | 2Gi | N/A |


## Variables


### Default Variables


| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_storage.namespace | Namespace for identity components | "tanzu-storage" | string | yes |
| tanzu_storage.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/{{ tanzu_storage.namespace }}" | string | yes |
| tanzu_storage.storage_class_defaults | Defaults for all storage classes when options don't explicitly override them in each storage class | See `common/vars/main.yaml` | dict | yes |
| tanzu_storage.ephemeral | Options related to ephemeral storage | See `common/vars/main.yaml` | dict | yes |
| tanzu_storage.aws | Options related to AWS storage | See `common/vars/main.yaml` | dict | yes |
| tanzu_storage.vmware | Options related to VMware storage | See `common/vars/main.yaml` | dict | yes |
| tanzu_storage.azure | Options related to Azure storage | See `common/vars/main.yaml` | dict | yes |
| tanzu_storage.pause_for_configuration | Pause to let use configure storage (in lieu of missing module). Only applicable to V7w/K8s when tanzu_cluster_is_v7wk8s is true | true | dict | yes |


### Additional Variables


| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_storage_classes_vmware | The storage classes for a VMware TKG cluster (see below for example) | - |  List of Dicts | When tanzu_provider == 'vmware' |
| tanzu_storage_classes_aws | The storage classes for an AWS TKG cluster (see below for example) | - | List of Dicts | When tanzu_provider == 'aws' |
| tanzu_storage_classes_azure | The storage classes for an Azure TKG cluster (see below for example) | - | List of Dicts | When tanzu_provider == 'azure' |

An example of storage classes configured for an AWS provider:

**NOTE:** only one default class is allowed.  If you specify a default class and one
already exists in your cluster, your class will result in `default=false` instead.

```yaml
tanzu_storage_classes_aws:
  - name:           ""
    default:        true
    type:           "gp2"   # io2, gp2, sc1, st1
    iops_per_gb:    "10"    # only applicable to io2
    file_system:    "ext4"
    encrypted:      "false"
    kms_key_id:     ""      # only applicable to encrypted volumes
    expandable:     false
    reclaim_policy: "Delete"
    mount_options:  []
```

An example of storage classes configured for a VMware provider:

**NOTE:** only one default class is allowed.  If you specify a default class and one
already exists in your cluster, your class will result in `default=false` instead.

```yaml
tanzu_storage_classes_vmware:
  - name:           "ssd_storage"
    datastores:
      - "datastore1"
    default:        false
    expandable:     false
    reclaim_policy: "Delete"
    file_system:    "ext4"
    mount_options:  []
    disk_format:    "thin"
```

An example of storage classes configured for an Azure provider:

**NOTE:** only one default class is allowed.  If you specify a default class and one
already exists in your cluster, your class will result in `default=false` instead.

```yaml
tanzu_storage_classes_azure:
  - name:           "azure-premium"
    default:        true
    type:           "Premium_LRS"    # azure storage account sku tier. only premium vms can attach Premium_LRS disks
    resource_group: ""               # resource group where azure disk will be created. default is k8s node resource group
    kind:           "dedicated"      # shared, dedicated, managed
    caching_mode:   "ReadWrite"      # ReadOnly, ReadWrite, None
    file_system:    "ext4"
    expandable:     false
    reclaim_policy: "Delete"
```

## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=storage make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=storage make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

**NOTE:** this is a work in progress.  See [Verification](#verification) for steps to manually demonstrate
the role.

In order to demonstrate the role:

```bash
ROLE=storage make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=storage make clean.role
```


## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
