# Security

* Deploy cert-manager for certificate management.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#security-hardening).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 500m | 1Gi| N/A - no persistent storage required |


## Variables


### Default Variables


| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_security.namespace | Namespace for security components | "tanzu-security" | string | yes |
| tanzu_security.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/{{ tanzu_security.namespace }}" | string | yes |
| tanzu_security.tls_providers | Providers to configure for creating CA Certs ("ca", "letsencrypt-stage", "letsencrypt-prod", "wildcard") | "ca" | string | yes |
| tanzu_security.default_resources | Default resource allocation. | Varies.  See `common/defaults/main.yaml` | dict | yes |
| tanzu_security.tls_root_ca_cert | Certificate for self signed root CA | "" | string | no |
| tanzu_security.tls_root_ca_key | Key for self signed root CA | "" | string | no |
| tanzu_security.actions.update_k8s_ca | Instructs RPK to update the Kube Nodes trusted CAs with the CA cert and reload Containerd | "true" | boolean | yes |

### tanzu_security.tls_providers

This variable determines the Cluster Issuer to set up for cert-manager to use valid settings are (`ca`, `letsencrypt-stage`, `letsencrypt-prod`, `wildcard`).

* `ca`: cert-manager will create self signed certs signed with a provided (see `tls_root_ca_cert` and `tls_root_ca_key`) or Generated CA keypair. The root CA Certificate will be placed on all Kubernetes nodes in `/etc/kubernetes/pki/tanzu_ca.crt`.
* `letsencrypt-stage`: cert-manager will request letsencrypt staging certificates. The letsencrypt staging CA Certificate will be placed on all Kubernetes nodes in `/etc/kubernetes/pki/tanzu_ca.crt`.
* `letsencrypt-prod`: cert-manager will request letsencrypt certificates. Letsencrypt has strict rate limits on requests of production certs, so this should /only/ be used on clusters that do not change frequently.
* `wildcard`: instead of requesting certs, it will use the provided certs.


#### tanzu_security.tls_root_ca_cert

This variable allows you to provide a plain text root ca cert that will be used to sign certificates when you set `tanzu_security.tls_provider ` to `ca`. If left empty, letsencrypt will create one for you. It's a good practice to provide a pre-created Root CA as often you will need to provide a Root CA Certificate to when installing Kubernetes (for API Server, and Docker/ContainerD) so that the Cluster will trust Certificates created with RPK ( OIDC Auth, Container Registry, etc).


#### tanzu_security.tls_root_ca_key

This variable allows you to provide a plain text root ca key that will be used to sign certificates when you set `tanzu_security.tls_provider ` to `ca`. If left empty, letsencrypt will create one for you. It's a good practice to provide a pre-created Root CA as often you will need to provide a Root CA Certificate to when installing Kubernetes (for API Server, and Docker/ContainerD) so that the Cluster will trust Certificates created with RPK ( OIDC Auth, Container Registry, etc).


#### tanzu_security.actions.update_k8s_ca

This is a semi-dangerous task that will attempt to update the Kubernetes nodes to trust the CA Certificate when using `ca` or `letsencrypt_stage` as your TLS Provider. It is tested to work on TKG 1.1 Clusters, and should be used with caution. Ideally you would perform this action out of band to RPK.


### Additional Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_security.ca.tls_root_ca_cert | Certificate for self signed root CA | "" | base64 | no |
| tanzu_security.ca.tls_root_ca_key | Key for self signed root CA | "" | base64 | no |
| tanzu_security.wildcard.tls_root_ca_cert | Certificate for CA that signed the wildcard | "" | base64 | no |
| tanzu_security.wildcard.tls_cert | Certificate for wildcard | "" | base64 | no |
| tanzu_security.wildcard.tls_key | Key for wildcard | "" | base64 | no |
| tanzu_security.letsencrypt_stage.tls_root_ca_cert | certificate for letsencrypt staging CA | "see common/main.yaml" | base64 | no |
| tanzu_security.actions.update_k8s_ca | Instructs RPK to update the Kube Nodes trusted CAs with the CA cert and reload Containerd | "true" | boolean | yes |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_ingress_domain | Subdomain used for deployment of RPK | - | string | yes |
| tanzu_ingress.class | Ingress class name (from ingress role) | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

If you are deploying RPK as a demo on a TKG Cluster, the defaults should be fine, however the certificates will be self signed, which may affect your ability to access images in the container registry.

If you wish to demonstrate RPK and have the cluster able to use images from the Container Registry it is highly recommended that you pre-create a CA key/cert and inject the cert into your Kubernetes cluster's Trusted CA Certs. You can then set `tanzu_security.tls_root_ca_cert` and `tanzu_security.tls_root_ca_key` with your CA key/cert pair and the Container Registry TLS certificates will be created using that CA key/cert pair and the Cluster should be able to pull from the Container Registry.

In production like situations you can set `tanzu_security.tls_provider` to `letsencrypt-prod` and all certificates will be signed with the trusted letsencrypt CA and your services should function with no further intervention.

> Currently there is a bug in the order that certificates and ingress are created that will require you to use a wildcard DNS entry pointing at your ingress IP if you use letsencrypt.

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=security make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=security make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

**NOTE:** this has not been implemented for this role yet.


## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=security make clean.role
```


## Author(s)

* [Dustin Scott](mailto:sdustin@vmware.com)
* [Paul Czarkowski](mailto:pczarkowski@vmware.com)
