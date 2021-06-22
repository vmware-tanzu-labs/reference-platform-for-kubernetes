# Kind Provider Documentation

## Table of Contents

- [Super Quick Start](#super-quick-start)
- [Detailed Quick Start](#detailed-quick-start)
   - [Ensure KIND Kubeconfig Cluster Contexts Exist](#ensure-kind-kubeconfig-cluster-contexts-exist)
   - [1. Create Cluster Inventory](#1-create-cluster-inventory)
   - [2. Set Variables](#2-set-variables)
     - [2a. Common Variables](#2a-common-variables)
       - [Common Variable Table](#common-variable-table)
     - [2b. Additional/Optional Variables (Expert)](#2b-additionaloptional-variables-expert)
       - [RPK Components](#rpk-components)
         - [Platform Profile](#platform-profile)
         - [Advanced Profile](#advanced-profile)
   - [3. Build RPK Container Image](#3-build-rpk-container-image)
   - [4. Deploy RPK to Your KIND Cluster](#4-deploy-rpk-to-your-kind-cluster)
     - [Deploy RPK](#deploy-rpk)

## Super Quick Start

1. `make setup.kind` - create the inventory file with global variables.

1. `make setup.kind.network` - deploy Calico as a CNI for KIND.

1. `kind get kubeconfig --name rpk >> ~/.kube/config` - ensure the KIND kubeconfig and contexts exist in the kubeconfig file.

1. `vim build/inventory.yaml` - edit and update the variables.

1. `make build` - build an image for deploying RPK.

1. `make deploy` - deploy RPK.

## Detailed Quick Start

This guide will take you through all of the necessary steps to prepare for and to deploy RPK.

### Ensure KIND Kubeconfig Cluster Contexts Exist

Make sure your local kubeconfig contains the reference to your deployed cluster:

```bash
kubectl config get-contexts
```

If the KIND kubeconfig context is missing, execute the following command:

```bash
kind get kubeconfig --name rpk >> ~/.kube/config
```

### 1. Create Cluster Inventory

A custom inventory will need to be created in order to define which cluster(s) you will be
deploying RPK to.  This inventory file is a [standard Ansible inventory file, in YAML
format](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#inventory-basics-formats-hosts-and-groups).  The inventory file additionally contains all of the necessary variables that feed into the deployment process.

To deploy the KIND cluster (tested with `v0.9.0`), inventory file, and CNI plugin, run the following commands:

```bash
make setup.kind
make setup.kind.network
```

> **NOTE:** The command above will drop your inventory file example at `build/inventory.yaml`.  If you need the inventory file at another location,
> you can run `INVENTORY=/path/to/my/inventory.yaml make setup.aws` for example.  Be sure that the `INVENTORY` variable matches this for following
> build and deploy steps.

### 2. Set Variables

After creating your inventory file, you will need to modify required variables that fit your environment needs:

```bash
vim build/inventory.yaml
```

#### 2a. Common Variables

##### Common Variable Table
Common variables, regardless of provider type, that will need to be modified are as follows:

| Variable Name         | Variable Type | Required | Default Value | Description                                                                              |
| --------------------- | ------------- | -------- | ------------- | ---------------------------------------------------------------------------------------- |
| rpk_profile           | string        | true     | (none)        | Name of the RPK profile (collection of cluster components) to deploy . Can be set to `platform` or `advanced`. See the contribution guide for more details on profiles [here](../CONTRIBUTING.md#profiles). |
| tanzu_cluster_name    | string        | true     | (none)        | The DNS friendly name (e.g. no underscores or special characters) of the cluster         |
| tanzu_default_tls_provider | string | true | "ca" | The default method for ingress TLS certificate deployment.  Can be set to `"ca"`, `"letsencrypt-stage"`, or `"letsencrypt-prod"`  See [Ingress TLS Certificates](../tls/tls.md) for expanded information on each setting. |
| tanzu_dns_provider    | string        | true     | route53       | The DNS provider/backend. Currently `route53`, `internal`, `nip.io`, and `local` are valid options.  See [Choosing the Appropriate DNS Solution](../DNS.md) for expanded information on each setting. |
| tanzu_ingress_domain  | string        | true     | (none)        | The ingress domain, which will be used to add ingress DNS records to                     |
| tanzu_kubectl_context | string        | true     | (none)        | The kubectl context (from kubeconfig) which represents the TKG target deployment cluster |

#### 2b. Additional/Optional Variables (Expert)

In most cases you will not need to manipulate any of the preset variables for RPK components.  However, you may encounter edge cases where the behavior of a component needs to be tuned or modified to affect how it is deployed into your cluster.  This is a rare circumstance, but if you would like to delve into all the possible configuration options for the deployment of RPK, see each of the component specific documentation pages from the table below.

##### RPK Components

###### Platform Profile

The following components are displayed in the order they will be deployed to a TKG cluster.  Additional tuning and behavior can be changed by adding additional variables to your `build/inventory.yaml`.

| Component Name | Documentation |
| --- | --- |
| workload-tenancy | [README.md](../../roles/components/core/workload-tenancy/README.md) |
| networking | [README.md](../../roles/components/core/networking/README.md) |
| security | [README.md](../../roles/components/core/security/README.md) |
| admission-control | [README.md](../../roles/components/core/admission-control/README.md) |
| storage | [README.md](../../roles/components/core/storage/README.md) |
| ingress | [README.md](../../roles/components/core/ingress/README.md) |
| identity | [README.md](../../roles/components/core/identity/README.md) |
| monitoring | [README.md](../../roles/components/core/monitoring/README.md) |
| container-registry | [README.md](../../roles/components/core/container-registry/README.md) |
| logging | [README.md](../../roles/components/core/logging/README.md) |
| autoscaling | [README.md](../../roles/components/core/autoscaling/README.md) |
| secret-management/hashicorp-vault | [README.md](../../roles/components/core/secret-management/hashicorp-vault/README.md) |
| application-stack | [README.md](../../roles/components/core/application-stack/README.md) |
| application-pipeline | [README.md](../../roles/components/core/application-pipeline/README.md) |

###### Advanced Profile

The following components are displayed in the order they will be deployed to a TKG cluster.  Additional tuning and behavior can be changed by adding additional variables to your `build/inventory.yaml`.

| Component Name | Documentation |
| --- | --- |
| networking | [README.md](../../roles/components/core/networking/README.md) |
| container-registry | [README.md](../../roles/components/core/container-registry/README.md) |
| tanzu-mission-control | [README.md](../../roles/components/extensions/tanzu-mission-control/README.md) |
| spring-cloud-gateway | [README.md](../../roles/components/extensions/spring-cloud-gateway/README.md) |
| spring-cloud-data-flow | [README.md](../../roles/components/extensions/spring-cloud-data-flow/README.md) |
| tanzu-application-catalog | [README.md](../../roles/components/extensions/tanzu-application-catalog/README.md) |
| tanzu-build-service | [README.md](../../roles/components/extensions/tanzu-build-service/README.md) |
| tanzu-observability | [README.md](../../roles/components/extensions/tanzu-observability/README.md) |
| tanzu-sql | [README.md](../../roles/components/extensions/tanzu-sql/README.md) |

### 3. Build RPK Container Image

Build a Docker image to run the Ansible playbooks.  This will relieve you of the need to install the Python dependencies locally.  To build the Docker image with the latest state of the project locally, run:

```bash
make build
```

To build this Docker image with custom image names (default: rpk) and image versions (default: latest), run the
following, being sure to substitute in your desired variables appropriately:

```bash
IMAGE=rpk-custom VERSION=v1.0.0 make build
```

> **NOTE:** *If using custom names, ensure that these match during the deploy stage.*

### 4. Deploy RPK to Your KIND Cluster

At this point, you should have completed all the necessary prerequisite work to deploy RPK to your workload cluster.  You are now ready to Deploy RPK.

#### Deploy RPK

You can deploy using the docker image with the latest state of the project locally.  To do so, run:

```bash
make deploy
```

To run when using custom image names, image versions, and/or inventory run the following, being sure to substitute your desired variables appropriately:

```bash
INVENTORY=/path/to/my/inventory.yaml IMAGE=rpk-test VERSION=v1.0.0 make deploy
```

> **NOTE:** *If you are deploying RPK to a Tanzu Mission Control provisioned cluster, you must run the following version of the command:*

```bash
TMC=true make deploy
```
