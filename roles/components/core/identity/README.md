# Identity

Deploy identity authentication using OIDC with the API server.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#identity-services).

## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 1 | 1.5Gi | N/A |

## Additional Requirements

If using the identity module with TKG1.3+, you will need to ensure that the native Dex/Gangway/Pinniped/LDAP/OIDC
integration is not setup during deployment time.  There is a conflict that is not supported at this time.

## Variables


### Default Variables


| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_identity.dex.dns | DNS name to use for Dex | oidc-server | string | yes |
| tanzu_identity.gangway.dns | DNS name to use for Gangway | oidc-client | string | yes |
| tanzu_identity.dex.tls_provider | TLS provider to use for Dex | {{ tanzu_security.tls_provider }} | string | yes |
| tanzu_identity.gangway.tls_provider | TLS provider to use for Gangway | {{ tanzu_security.tls_provider }} | string | yes |
| tanzu_identity.namespace | Namespace for identity components | "tanzu-identity" | string | yes |
| tanzu_identity.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/{{ tanzu_identity.namespace }}" | string | yes |
| tanzu_identity.{dex,gangway,ldap}.resources | Default resource allocation for identity components | Varies.  See `common/vars/main.yaml` | dict | yes |


### Additional Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_ingress_domain | Subdomain used for deployment of RPK | - | string | yes |
| ingress_ip | Ingress IP Address (from ingress role) | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* security
* ingress


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=identity make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=identity make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.  This demo simply
verifies connectivity to the OIDC client, prints the URL to visit to authenticate the client,
creates ClusterRoleBindings for some sample users,
and prints out authentication info of the sample users to test the client with:

```bash
ROLE=identity make demo.role
```

Sample output:

```bash
...
ok: [vmware_tanzu_cluster1] => (item={'username': 'john', 'password': 'bar', 'first_name': 'john', 'last_name': 'doe', 'email': 'john@vsphere.ihazcloudthings.xyz', 'clusterrole': 'view', 'clusterrolebinding': 'john@vsphere.ihazcloudthings.xyz'}) => {
    "msg": [
        "User: john@vsphere.ihazcloudthings.xyz",
        "Password: bar",
        "You can access the OIDC Client Dashboard at URL https://oidc-client.vsphere.ihazcloudthings.xyz",
        "You will need to manually trust the CA certificates for this to work"
    ]
}
```

Once the role has been demonstrated, you can test it by visiting the URL printed in the
demo message:

```bash
        "You can access the OIDC Client Dashboard at URL https://oidc-client.vsphere.ihazcloudthings.xyz"
```

You can authenticate with one of the users from the demo message:

```bash
        "User: john@vsphere.ihazcloudthings.xyz",
        "Password: bar",
```

Once you have successfully authenticated, you will be able to download your kubeconfig to access
your cluster.  Use your kubeconfig following the download by running:

```bash
export KUBECONFIG=/path/to/downloaded/kubeconf.txt
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=identity make clean.role
```


## Author(s)
[Rich Lander](mailto:landerr@vmware.com)
