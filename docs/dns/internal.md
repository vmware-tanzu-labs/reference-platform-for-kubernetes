# Internal DNS Configuration

1. Within your `build/inventory.yaml` file, which was created with `make setup.<provider>` set the variable `tanzu_dns_provider: "internal"`.

1. Within your `build/inventory.yaml` file, set the variable `tanzu_ingress_domain` to a domain name of your choosing.  Since this configuration is typically used for testing, the choice here is non-consequential.  Example: `tanzu_ingress_domain: "my.fake.domain.com"`

1. Continue setting up your `build/inventory.yaml`: [Quick Start: Provider Specific Steps](../QUICKSTART.md#provider-specific-steps)

1. [Deploy RPK](../QUICKSTART.md#deploy)

1. Gather the internal DNS IP address:
```bash
kubectl get service -n tanzu-dns
```
You may want to obtain the IP address in cases where the `EXTERNAL IP` field is presented as a fully-qualified domain name:
```bash
nslookup a.b.c.d.com
```

6. Set your local workstation's DNS resolver (DNS servers) to the value from the `EXTERNAL IP` field in the previous command.    See [Setting DNS Resolvers](setting-dns-resolvers.md) for more details.

1. You may now access your ingress resources via names.

1. For further validation that the internal DNS service is functioning, please see the validation steps [here](../../roles/components/core/ingress/README.md#verification)

[Back to Quick Start](../QUICKSTART.md#dns-options)
