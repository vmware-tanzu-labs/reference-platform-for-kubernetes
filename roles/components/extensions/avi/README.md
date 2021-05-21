# Avi

This module installs the Avi Kubernetes Operator into a single cluster.

**NOTE:** a working Avi Vantage installation in vCenter is a prerequisite for this module to work.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md)

## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 500m | 1Gi | N/A - no persistent storage required |


## Variables

### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| avi.namespace | Namespace for Avi components | "avi" | string | yes |
| avi.staging_dir | Local directory to write the staging manfiests to | "/tmp/staging/extensions/avi" | string | yes |
| avi.vcenter.folder | vCenter folder where Avi VMs are to be deployed | "Discovered Virtual Machines" | string | No |
| avi.controller.password | Default password for the Avi "admin" user (also used as the backup passphrase) | "tanzu" | string | No |
| avi.controller.ha_mode | Whether to deploy a 3-node HA Cluster or not | false | boolean | No |
| avi.controller.vip | Cluster VIP for Avi Cluster | "" | string | Only for `ROLE=avi make infra` and `avi_controller_ha_mode == true` |
| avi.controller.memory_mb | Memory, in megabytes, to set for Avi Controller(s) | "24576" | string | Only for `ROLE=avi make infra` |
| avi.controller.num_cpus | Number of CPUs to set for Avi Controller(s) | "8" | string | Only for `ROLE=avi make infra` |
| avi.controller.disk_size_gb | Size, in gigabytes, to set for Avi Controller(s) disk | "128" | string | Only for `ROLE=avi make infra` |
| avi.controller.anti_affinity_create | Whether to create anti-affinity rules for Avi VMs | false | boolean | Only for `ROLE=avi make infra` and `avi_controller_ha_mode == true` |
| avi.service_engine.memory_mb | Memory, in megabytes, to set for Avi Service Engine(s) | "4096" | string | Only for `ROLE=avi make infra` |
| avi.service_engine.num_cpus | Number of CPUs to set for Avi Service Engine(s) | "2" | string | Only for `ROLE=avi make infra` |
| avi.service_engine.disk_size_gb | Size, in gigabytes, to set for Avi Service Engine(s) disk | "20" | string | Only for `ROLE=avi make infra` |

### Additional Variables

The following variables must be set for proper operation of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| avi_vcenter_datacenter | vCenter datacenter where Avi VMs are to be deployed | "" | string | Only for `ROLE=avi make infra` |
| avi_vcenter_cluster | vCenter cluster where Avi VMs are to be deployed | "" | string | Only for `ROLE=avi make infra` |
| avi_vcenter_datastore | vCenter datastore where Avi VMs are to be deployed | "" | string | Only for `ROLE=avi make infra` |
| avi_vcenter_management_portgroup | vCenter portgroup for Avi management traffic | "" | string | Only for `ROLE=avi make infra` |
| avi_controller_dns_server | DNS server to set for the Avi controller | "" | string | Only for `ROLE=avi make infra` |
| avi_controller_ntp_server | NTP server to set for the Avi controller | "" | string | Only for `ROLE=avi make infra` |
| avi_controller_networks | Avi Networks for all Avi nodes | [] | List of Dicts (see [sample](#sample-variables) below) | Only for `ROLE=avi make infra` |
| avi_controller_license_data | The license data for Avi (see [sample](#sample-variables) below) | "" | string | Only for `ROLE=avi make infra` |
| avi_service_engine_networks | Avi Service Engine Data Networks for all Avi SEs | [] | List of Dicts (see [sample](#sample-variables) below) | Only for `ROLE=avi make infra` |
| avi_service_engine_management_ip | Management IP for communications back to Avi Controller | "" | string | Only for `ROLE=avi make infra` |
| avi_service_engine_management_netmask | Management Netmask for communications back to Avi Controller | "" | string | Only for `ROLE=avi make infra` |
| avi_service_engine_management_gateway | Management Gateway for communications back to Avi Controller | "" | string | Only for `ROLE=avi make infra` |
| tanzu_load_balancer_subnet | The subnet in which to alloate IP addresses from in format `1.2.3.4/24` | "" | string | yes |
| tanzu_load_balancer_start_ip | The first IP address in a range of IP addresses to allocate for load balancers | "" | string | yes |
| tanzu_load_balancer_end_ip | The last IP address in a range of IP addresses to allocate for load balancers | "" | string | yes |

#### Sample Variables

Below is a set of sample variables as provided to the `build/inventory.yaml` file prior to running [deployment](#deploying):

```yaml
    # vcenter vars
    avi_vcenter_datacenter:           "My-Datacenter"
    avi_vcenter_cluster:              "My-Cluster01"
    avi_vcenter_datastore:            "datastore01"
    avi_vcenter_folder:               "My-Datacenter/vm/Virtual Machines/Linux"
    avi_vcenter_management_portgroup: "VM Network"

    # avi controller vars
    avi_controller_password:      "tanzu"
    avi_controller_ha_mode:       false
    avi_controller_dns_server:    "8.8.8.8"
    avi_controller_ntp_server:    "0.pool.ntp.org"
    avi_controller_networks:
      - ip:      "172.26.0.6"
        netmask: "24"
        gateway: "172.26.0.1"
        fqdn:    "avi1.example.com"
    avi_controller_anti_affinity_create: false
    avi_controller_license_data: |
      avi_license data here

    # service engine vars
    avi_service_engine_networks:
      - portgroup: "VM Network"
        ip:        "172.26.0.9"
        netmask:   "24"
    avi_service_engine_management_ip:      "172.26.0.10"
    avi_service_engine_management_netmask: "24"
    avi_service_engine_management_gateway: "172.26.0.1"
```


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* `extensions/avi/infra`

**NOTE:** the above is a module that gets deployed out of band.  It MUST be run before this role will work.  See [Deploying](#deploying)
below.


## Deploying

Prior to deployment, the vCenter OVA for Avi must be downloaded.  This OVA will be set as the `AVI_OVA` environment variable
when running the deployment command.  You may obtain the Avi control plane OVA from
https://vault.vmware.com/group/nsx/avi-networks-technical-resources.

**NOTE:** roles from `.dependencies.yaml` are also deployed.

Prior to deployment, the infrastructure in vCenter must be initialized.  To do so, be sure to fill out all
required variables from [Variables](#additional-variables) and then run the following command:

```bash
AVI_OVA=/path/to/downloaded/avi/ova-file.ova ROLE=extensions/avi make infra
```

In order to deploy the role from a versioned image:

```bash
ROLE=extensions/avi make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=extensions/avi make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

Brief description of what the demonstration does here.

In order to demonstrate the role:

```bash
ROLE=extensions/avi make demo.role
```

Sample output:

```bash
...
relevant sample output here
```

Accessing the service:

```bash
relevant sample output here (username, password, url, ip, etc)
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

**NOTE:** this does not cleanup the infrastructure, only the Avi Kubernetes Operator (AKO).

```bash
ROLE=extensions/avi make clean.role
```


## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
