# AWS Route53 DNS Instructions

## Prerequisites
First, create the public subdomain in AWS Route 53.  Once the subdomain is created, note the nameservers which are serving this domain:

![Custom Subdomain](../images/custom-subdomain.png)

In order to delegate your subdomain, you will need to create the nameserver DNS entries (captured in previous step) in your DNS provider.  This is **required** until we can solve
for a more consistent DNS experience and potentially support more DNS providers.  The variables that tie to your DNS hostnames should also be updated in your inventory (see below).

| Use | Record Type | Example Host | Example Record |Associated Variable |
| --- | --- | --- | --- | --- |
| Subdomain (created in example.net domain DNS servers) | NS Record | subdomain.example.net | ns-218.awsdns-27.com (see image above)  | tanzu_ingress_domain |

## Setup the Inventory File

1. Within your `build/inventory.yaml` file, which was created with `make setup.<provider>` set the variable `tanzu_dns_provider: "route53"`.

2. Within your `build/inventory.yaml` file, set the variable `tanzu_ingress_domain` to the domain you configured in the route53 service. From the example in the [Prerequisites](#prerequisites) section, this would be set like this: `tanzu_ingress_domain: "subdomain.example.net"`

3. Continue setting up your `build/inventory.yaml`: [Quick Start: Provider Specific Steps](../QUICKSTART.md#provider-specific-steps)

4. [Deploy RPK](../QUICKSTART.md#deploy)
