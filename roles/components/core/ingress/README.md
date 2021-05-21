# Ingress

This role deploys Contour and Envoy components for Ingress as well as ExternalDNS for external cluster DNS management.  In vSphere environments, a Metal LB load balancer
to provide Service Type of Load Balancer is also deployed.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#service-routing).

## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 2 | 4Gi | N/A - no persistent storage required |

**NOTE:** above sizing requirements will adjust based on cluster size due to
components utilizing a DaemonSet.

## Variables

### Default Variables

The following variables are defaulted in the `common/defaults/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_ingress.namespace | Namespace for tanzu_ingress components | "tanzu-ingress" | string | yes |
| tanzu_ingress.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-ingress" | string | yes |
| tanzu_ingress.default_resources | Normal Kubernetes resource construct defining default resources for components when undefined | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.contour.image | Image to use for contour components | "quay.io/bitnami/contour" | string | yes |
| tanzu_ingress.contour.image_tag | Image tag to use for contour components | "1.8.1" | string | yes |
| tanzu_ingress.contour.resources | Normal Kubernetes resource construct defining resource requirements | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.envoy.image | Image to use for envoy components | "quay.io/bitnami/envoy" | string | yes |
| tanzu_ingress.envoy.image_tag | Image tag to use for envoy components | "1.15.0" | string | yes |
| tanzu_ingress.envoy.resources | Normal Kubernetes resource construct defining resource requirements | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.metallb.controller.image | Image to use for Metal LB controller components | "metallb/controller" | string | yes |
| tanzu_ingress.metallb.controller.image_tag | Image tag to use for Metal LB controller components | "v0.8.1" | string | yes |
| tanzu_ingress.metallb.speaker.image | Image to use for Metal LB speaker components | "metallb/speaker" | string | yes |
| tanzu_ingress.metallb.speaker.image_tag | Image tag to use for Metal LB speaker components | "v0.8.1" | string | yes |
| tanzu_ingress.metallb.resources | Normal Kubernetes resource construct defining resource requirements | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.external_dns.image | Image to use for external-dns components | "k8s.gcr.io/external-dns/external-dns" | string | yes |
| tanzu_ingress.external_dns.image_tag | Image tag to use for external-dns components | "v0.7.3" | string | yes |
| tanzu_ingress.bind.image | Image to use for internal DNS configuration (bind9) | "internetsystemsconsortium/bind9" | string | yes |
| tanzu_ingress.bind.image_tag | Image tag to use for internal DNS configuration (bind9) | "9.16" | string | yes |
| tanzu_ingress.bind_init.image | Image to use for initializing internal DNS configuration (bind9) | "alpine" | string | yes |
| tanzu_ingress.bind_init.image_tag | Image tag to use for initializing internal DNS configuration (bind9) | "3.12.1" | string | yes |
| tanzu_ingress.external_dns.resources | Normal Kubernetes resource construct defining resource requirements | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.bind.resources | Normal Kubernetes resource construct defining resource requirements for internal DNS configuration (bind9) | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.bind_init.resources | Normal Kubernetes resource construct defining resource requirements for internal DNS configuration (bind9) | See `common/defaults/main.yaml` | dict | yes |
| tanzu_ingress.external_dns.known_providers | Known configurable providers for external-dns.  *NOTE*: Future providers may be added in the future. | ["route53", "internal", "azure"] | list | yes |
| tanzu_ingress.external_dns.delete_records | Whether to allow external-dns to delete records or not.  If set to `false` external-dns will only add and update records. | true | boolean | yes |
| tanzu_ingress.external_dns.sync_interval_seconds | The interval in which external-dns will check to see if it needs to update records. | 15 | integer | yes |

### Additional Variables

The following variables must be set for proper operatiion of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_load_balancer_start_ip | The first IP address in a range of IP addresses to allocate for load balancers | - | string | yes |
| tanzu_load_balancer_end_ip | The last IP address in a range of IP addresses to allocate for load balancers | - | string | yes |
| tanzu_ingress_provider | Backend DNS service to configure external-dns's components for. Available providers are `route53`, `internal` | "route53" | string | yes |
| aws_access_key | AWS Access key used to update DNS records (`route53` backend only) | - | string | yes |
| aws_secret_key | AWS Secret key used to update DNS records (`route53` backend only) | - | string | yes |
| tanzu_ingress_aws_eip_allocation_id | EIP Allocation ID of a reserved IP to be used for Ingress (example `eipalloc-0bd4399a62318d2bb`) | - | string | no |
| azure_tenant_id | Azure Tenant ID used to update DNS records in Azure DNS Zone (`azure` backend only) | - | string | yes |
| azure_subscription_id | Azure Subscription ID used to update DNS records in Azure DNS Zone (`azure` backend only) | - | string | yes |
| azure_client_id | Azure Client ID used to update DNS records in Azure DNS Zone (`azure` backend only) | - | string | yes |
| azure_client_secret | Azure Client secret used to update DNS records in Azure DNS Zone (`azure` backend only) | - | string | yes |
| azure_dns_resource_group_name | Name of the Azure Resource Group that contains the configured Azure DNS zone (`azure` backend only) | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* security


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=ingress make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=ingress make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

**NOTE:** this is a work in progress.  See [Verification](#verification) for steps to manually demonstrate
the role.

```bash
ROLE=ingress make demo.role
```


## Verification

### Contour

Follow steps located at https://projectcontour.io/getting-started/#example-workload for verification of ingress.

### AWS Route53 DNS Backend
See:
* https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#verify-externaldns-works-ingress-example
* https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#verify-externaldns-works-service-example
to verify proper operation of external-dns.

### Internal DNS Backend

```bash
# Observe the running external-dns pod and note the pod name
kubectl get pods -n tanzu-dns

# Check the container logs
## You should observe the addition of records in both sets of logs
kubectl logs [pod name] -c bind -n tanzu-dns
kubectl logs [pod name] -c external-dns -n tanzu-dns

# Obtain the service IP Address/Hostname and note the EXTERNAL-IP
kubectl get service -n tanzu-dns

# Test that you can look up an ingress or service
nslookup registry.[your domain] [EXTERNAL IP]
```

## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=ingress make clean.role
```


## Author(s)
* [James Kirkland](mailto:kirklandja@vmware.com)
* [Paul Czarkowski](mailto:pczarkowski@vmware.com)
* [Andrew J. Huffman](mailto:ahuffman@vmware.com)
* [Paul Wiggett](mailto:pwiggett@vmware.com)
