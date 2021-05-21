# RPK Components

## Platform Profile

The following components are displayed in the order they will be deployed to a TKG cluster.  Additional tuning and behavior can be changed by adding additional variables to your `build/inventory.yaml`.

| Component Name | Documentation |
| --- | --- |
| workload-tenancy | [README.md](../roles/components/core/workload-tenancy/README.md) |
| networking | [README.md](../roles/components/core/networking/README.md) |
| security | [README.md](../roles/components/core/security/README.md) |
| admission-control | [README.md](../roles/components/core/admission-control/README.md) |
| storage | [README.md](../roles/components/core/storage/README.md) |
| ingress | [README.md](../roles/components/core/ingress/README.md) |
| identity | [README.md](../roles/components/core/identity/README.md) |
| monitoring | [README.md](../roles/components/core/monitoring/README.md) |
| container-registry | [README.md](../roles/components/core/container-registry/README.md) |
| logging | [README.md](../roles/components/core/logging/README.md) |
| secret-management/hashicorp-vault | [README.md](../roles/components/core/secret-management/hashicorp-vault/README.md) |
| secret-management/etcd-encryption | [README.md](../roles/components/core/secret-management/etcd-encryption/README.md) |
| application-stack | [README.md](../roles/components/core/application-stack/README.md) |
| application-pipeline | [README.md](../roles/components/core/application-pipeline/README.md) |
| service-mesh/istio | [README.md](../roles/components/core/service-mesh/istio/README.md) |

## Advanced Profile

The following components are displayed in the order they will be deployed to a TKG cluster.  Additional tuning and behavior can be changed by adding additional variables to your `build/inventory.yaml`.

| Component Name | Documentation |
| --- | --- |
| networking | [README.md](../roles/components/core/networking/README.md) |
| container-registry | [README.md](../roles/components/core/container-registry/README.md) |
| extensions/tanzu-mission-control | [README.md](../roles/components/extensions/tanzu-mission-control/README.md) |
| extensions/avi | [README.md](../roles/components/extensions/avi/README.md) |
| extensions/spring-cloud-gateway | [README.md](../roles/components/extensions/spring-cloud-gateway/README.md) |
| extensions/spring-cloud-data-flow | [README.md](../roles/components/extensions/spring-cloud-data-flow/README.md) |
| extensions/tanzu-application-catalog | [README.md](../roles/components/extensions/tanzu-application-catalog/README.md) |
| extensions/tanzu-build-service | [README.md](../roles/components/extensions/tanzu-build-service/README.md) |
| extensions/tanzu-observability | [README.md](../roles/components/extensions/tanzu-observability/README.md) |
| extensions/tanzu-service-mesh | [README.md](../roles/components/extensions/tanzu-service-mesh/README.md) |
| extensions/tanzu-sql | [README.md](../roles/components/extensions/tanzu-sql/README.md) |
