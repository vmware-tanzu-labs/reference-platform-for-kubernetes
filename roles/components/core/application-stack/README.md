# Application Stack

## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md).

This [Guide](https://tanzu.vmware.com/developer/guides/spring/) demonstrates design, development, and
deployment of Spring Boot microservices on Kubernetes using the following concepts:

* Externalized configuration using ConfigMaps, Secrets, and PropertySource
* Kubernetes API server access using ServiceAccounts, Roles, and RoleBindings
* Health checks using Application Probes
  * readinessProbe
  * livenessProbe
  * startupProbe
* Reporting application state using Spring Boot Actuators
* Service discovery across namespaces using DiscoveryClient
* Exposing API documentation using Swagger UI
* Building a Docker image using best practices
* Layering JARs using the Spring Boot plugin
* Observing the application using Prometheus exporters
* Secure secret management using Hashicorp Vault
* Horizontal Pod auto-scaling


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| 8 | 8Gi | *100Gi |

*= Large storage requirement due to Prometheus monitoring dependency


## Variables


### Default Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_app_stack.{department,employee,organization,monitoring,mongodb,gateway}.namespace | Namespace name for microservices components | "tanzu-app-stack-{}" | string | yes |
| tanzu_app_stack.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-app-stack" | string | yes |
| tanzu_app_stack.service_account | Service account to use for microservices | "tanzu-app-stack" | string | yes |
| tanzu_app_stack.resource_name | Default resource name form all resources | "tanzu-app-stack" | string | yes |
| tanzu_app_stack.build_images | Build Docker images locally | false | bool | yes |
| tanzu_app_stack.push_images | Push locally built Docker images to remote image repository | false | bool | yes |
| tanzu_app_stack.remote_image_repo | Repo for push locally built Docker images to remote image repository | false | bool | yes |
| tanzu_app_stack.remote_image_username | Username for push locally built Docker images to remote image repository | false | bool | yes |
| tanzu_app_stack.remote_image_password | Password for push locally built Docker images to remote image repository | false | bool | yes |
| tanzu_app_stack.remote_image_repo | Remote image repository to push images to | github.io | string | yes |
| tanzu_app_stack.remote_image_username | Remote image repository user, needs to be replaced with real user name | andriykalashnykov | string | yes |
| tanzu_app_stack.remote_image_password | Remote image repository password | _replace_me_with_real_pwd_ | string | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.remote_image | Image to pull | "andriykalashnykov/{}-debug" | string | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.local_image | Image to push | "andriykalashnykov/{}-debug" | string | yes |
| tanzu_app_stack.mongodb.image | Image to pull for Mongo | "mongo" | string | yes |
| tanzu_app_stack.mongodb.image_tag | Image tag to pull for Mongo | "4.2.3" | string | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.remote_image_tag | Image tag to pull | "1.2" | string | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.local_image_tag | Image tag to push | "1.2" | string | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.git_path | Path to Dockerfile in remote Git Repo (only for push/pull) | "{}-servuce" | string | yes |
| tanzu_app_stack.{department,employee,organization,monitoring,mongodb,gateway}.resource_name | Default resource name for microservices components | "tanzu-app-stack-{}" | string | yes |
| tanzu_app_stack.{department,employee,organization,monitoring,mongodb,gateway}.resources | Default resource allocation microservices components | Varies.  See `common/vars/main.yaml` | dict | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.replicas | Number of replicas to deploy for each microservice (can be scaled by HPA) | 1 | integer | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.min_replicas | Minimum number of replicas to keep (HPA) | 1 | integer | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.max_replicas | Maximum number of replicas to deploy (HPA) | 3 | integer | yes |
| tanzu_app_stack.{department,employee,organization,gateway}.target_utilization | Target utilization before scaling (HPA) | "50" | string | yes |
| tanzu_app_stack.gateway.dns | DNS for gateway | "company.{{ tanzu_ingress_domain }}" | string | yes |
| tanzu_monitoring.gateway.app_prefix | Application path to access the stack (e.g. http://hostname/{{ tanzu_monitoring.gateway.app_prefix }}) | "/tanzu-app-stack/gateway" | string | yes |
| tanzu_app_stack.monitoring.enabled | Add Promethues monitoring services. Requires Prometheus | true | bool | yes |
| tanzu_app_stack.monitoring.dns | DNS entry to access the prometheus instance. | "company-monitor.{{ tanzu_ingress_domain }}" | string | yes |
| tanzu_app_stack.monitoring.metrics_path | Metrics path to access the prometheus instance. | "/actuator/prometheus" | string | yes |
| tanzu_app_stack.vault.enabled | Use Hashicorp Vault to store Mongo DB login/password. Requires Vault. If value is `false` application stack will use mounted secret instead | true | bool | yes |


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
* monitoring (only when `tanzu_app_stack.monitoring.enabled` is `true`)
* secret-management/hashicorp-vault (only when `tanzu_app_stack.vault.enabled` is `true`)


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=application-stack make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=application-stack make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

```bash
ROLE=application-stack make demo.role
```

This will create three jobs in the `tanzu-app-stack-gateway` namespace that
load dummy data into the app stack database and then apply traffic to the app
by querying that data.  Traffic will be sent to the app stack for 5
minutes.  You can watch the CPU percentages change and the gateway deployment
scale as the traffic hits the app witt the following command:

```bash
watch kubectl get hpa -A
```

Default variables for demo:

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_app_stack.demo.traffic_app_name | Application name for traffic demo | tanzu-app-stack-traffic | string | yes |
| tanzu_app_stack.demo.traffic_rate | Requests per second sent to the app sttack by each traffic job | 50 | string | yes |
| tanzu_app_stack.demo.traffic_duration | Duration that traffic jobs run | 300s | string | yes |


### Accesing the Demo App

#### Swagger UI

http://{{ tanzu_app_stack.gateway.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/swagger-ui.html

#### REST API

http://{{ tanzu_app_stack.gateway.dns }}/department
http://{{ tanzu_app_stack.gateway.dns }}/employee
http://{{ tanzu_app_stack.gateway.dns }}/organization

#### Prometheus metrics

http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/department/actuator/metrics
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/employee/actuator/metrics
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/organization/actuator/metrics
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/actuator/metrics

http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/department/actuator/prometheus
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/employee/actuator/prometheus
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/organization/actuator/prometheus
http://{{ tanzu_app_stack.monitoring.dns }}/{{ tanzu_app_stack.gateway.app_prefix }}/actuator/prometheus


## Cleaning

To remove the role, from the root of the repo:

**NOTE:** only this role is removed and not the role dependencies.

```bash
ROLE=application-stack make clean.role
```


## Author(s)

[Andriy Kalashnykov](mailto:akalashnykov@vmware.com)
[Rich Lander](mailto:landerr@vmware.com)
