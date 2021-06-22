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
