# Logging

Deploy the EFK stack.  See https://tanzu.vmware.com/developer/guides/kubernetes/observability/ for more information.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#observability).

## Variables

### Default Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_logging.namespace | Namespace name for logging components | "tanzu-logging" | string | yes |
| tanzu_logging.dns | DNS name to use for Kibana | {{ tanzu_logging }}.{{ tanzu_ingress_domain }} | string | yes |
| tanzu_logging.elastic.password | Password for the elastic user for kibana/elastic | "tanzu" | string | yes |
| tanzu_logging.elastic.cluster_name | Cluster identifier for log record filtering | "cluster-0" | string | yes |
| tanzu_logging.elastic.instance_name | Instance identifier for log record filtering | "instance-0" | string | yes |
| tanzu_logging.operator.image | Container image for elasticsearch operator | "docker.elastic.co/eck/eck-operator" | string | yes |
| tanzu_logging.operator.image_version | Container image version for elasticsearch operator | "1.1.2" | string | yes |
| tanzu_logging.elastic.version | Elasticsearch version to use in elasticsearch manifest | "7.8.0" | string | yes |
| tanzu_logging.fluent.image | Container image for fluentbit | "registry.tkg.vmware.run/fluent-bit" | string | yes |
| tanzu_logging.fluent.image_version | Container image for fluentbit | "v1.5.3_vmware.1" | string | yes |
| tanzu_logging.kibana.version | Kibana version to use in kibana manifest | "7.8.0" | string | yes |
| tanzu_logging.{elastic,operator,fluent,kibana}.resources | Default resource allocation microservices components | Varies.  See `common/defaults/main.yaml` | dict | yes |
| tanzu_logging.tls.secret | secret name for TLS | "tanzu-logging-tls" | string | yes |
| tanzu_logging.tls.provider | TLS provider | "{{ tanzu_default_tls_provider }}" | string | yes |

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
ROLE=logging make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=logging make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=logging make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=logging make clean.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=logging make clean.test.role
```


## Accessing the Dashboard

After the role has been deployed, you can access your dashboard with the following information:

* **URL:** `https://{{ tanzu_logging.dns }}.{{ tanzu_ingress_domain }}` (e.g. logs.mydomain.com)
* **Username:** `{{ tanzu_logging.elastic.username }}` (default: elastic)
* **Password:** `{{ tanzu_logging.elastic.password }}` (default: tanzu)


## Author(s)
[Matyas Danter](mailto:mdanter@vmware.com)
