# Monitoring

Set up monitoring and alerting for a TKG cluster.

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#cluster-monitoring).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 4 | 8Gi | 100Gi |

**NOTE:** above sizing requirements will adjust based on cluster size due to
components utilizing a DaemonSet.


## Variables


### Default Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_monitoring.namespace | Namespace name for monitoring components | "tanzu-monitoring" | string | yes |
| tanzu_monitoring.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/{{ tanzu_identity.namespace }}" | string | yes |
| tanzu_monitoring.{prometheus,grafana,alertmanager}.dns | DNS to use for each item | "http://{{ item }}" | string | yes |
| tanzu_monitoring.{prometheus,alertmanager,grafana,node_exporter,rbac_proxy,kube_state_metrics,operator,adapter}.{image,image_tag,resources} | Container image information for components | See `common/vars` | string | yes |
| tanzu_monitoring.alertmanager.slack_url | Slack URL to send default alerts to | none | string | no |
| tanzu_monitoring.alertmanager.watchdog_slack_url | Slack URL to send watchdog alerts to | none | string | no |
| tanzu_monitoring.alertmanager.critical_slack_url | Slack URL to send critical alerts to | none | string | no |

### Additional Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |
| tanzu_ingress_domain | Subdomain used for deployment of rpk | - | string | yes |
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
ROLE=monitoring make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=monitoring make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

**NOTE:** this is a work in progress.  See [Verification](#verification) for steps to manually demonstrate
the role.

```bash
ROLE=monitoring make demo.role
```


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=monitoring make clean.role
```


## Author(s)
[Rich Lander](mailto:landerr@vmware.com)
