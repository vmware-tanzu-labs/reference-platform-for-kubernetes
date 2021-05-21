# NetworkManager Instructions

1. Within your `build/inventory.yaml` file, which was created with `make setup.<provider>` set the variable `tanzu_dns_provider: "local"`.

1. Within your `build/inventory.yaml` file, set the variable `tanzu_ingress_domain` to a domain name of your choosing.  Since this configuration is typically used for testing, the choice here is non-consequential.  Example: `tanzu_ingress_domain: "my.fake.domain.com"`

1. Continue setting up your `build/inventory.yaml`: [Quick Start: Provider Specific Steps](../QUICKSTART.md#provider-specific-steps)

1. [Deploy RPK](../QUICKSTART.md#deploy)

## Built in NetworkManager / DNSmasq

Ensure that NetworkManager is using dnsmasq

/etc/NetworkManager/conf.d/00-use-dnsmasq.conf:
```
[main]
dns=dnsmasq
```

Set additional hosts to include a `/etc/hosts.rpk`:

/etc/NetworkManager/dnsmasq.d/00-home.conf:
```
addn-hosts=/etc/hosts
addn-hosts=/etc/hosts.rpk
```

After you have deployed RPK, copy `build/hosts` to `/etc/hosts.rpk`

```bash
sudo cp build/hosts /etc/hosts.rpk
```

Restart network manager:

```bash
systemctl restart NetworkManager.service
```

[Back to Quick Start](../QUICKSTART.md#dns-options)
