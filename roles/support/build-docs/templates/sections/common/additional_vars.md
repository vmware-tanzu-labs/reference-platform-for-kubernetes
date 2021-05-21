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
