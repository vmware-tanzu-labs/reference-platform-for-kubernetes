# nip.io Instructions

> Some DNS servers / routers have **Rebind protection - Discard upstream RFC1918 responses** which will block nip.io and similar domains from resolving local IP addresses. You may need to check your router (or dnsmasq) configuration to disable that setting if you can't resolve addresses for your local network.

This is the easiest way to solve it, but it does rely on an external service, so might be an issue in a firewalled environment.

If you don't know the IP address that will be assigned to your ingress service, you'll need to run the ingress role first.


## KIND

1. Download and install Kubernetes IN Docker 0.9.0 from [here](https://github.com/kubernetes-sigs/kind/releases/tag/v0.9.0).

1. Configure and set up KIND cluster (the KIND IP should be your machines primary IP address, the following will work on linux):

  ```bash
  make setup.kind
  ```

**NOTE:** the above assumes the Docker flag `--network=host`.  If this will not work for some reason, you will need to set
the real IP of your system in the `build/inventory.yaml` (created by `make setup.<provider>`) file as the `ingress_ip` variable.
You can get the real IP of your system using something similar to `hostname -I | awk '{ print $1}'`.

3. Deploy Calico as KIND CNI
  ```bash
  make setup.kind.network
  ```
4. Within your `build/inventory.yaml` file, which was created with `make setup.<provider>` set the variable `tanzu_dns_provider: "local"`.

1. Continue setting up your `build/inventory.yaml`: [Quick Start: Provider Specific Steps](../QUICKSTART.md#provider-specific-steps)

1. [Deploy RPK](../QUICKSTART.md#deploy)


## TMC / AWS

1. Create an inventory for your cluster:

  ```bash
  make setup.aws
  ```

2. Remove or comment out (e.g `#aws_access_key: ""`) `aws_access_key` and `aws_secret_key` from your `build/inventory.yaml`.

1. Update any other `build/inventory.yaml` values for your cluster.

1. Deploy the `ingress` role (we need to do this first so we can collect our Ingress IP address):
  ```bash
  ROLE=ingress make deploy.test.role
  ```

5. Run the `common/ingress-ip` role which will politely return the IP address needed:

  ```bash
  $ ROLE=common/ingress-ip make deploy.test.role
  ...
  ...
  TASK [common/ingress-ip : inspect the ingress ip] ******************************
  ok: [aws_tanzu_cluster1] =>
  ingress_ip: 3.18.48.175
  ```

6. Set `tanzu_ingress_domain=3.18.48.175.nip.io` in your `build/inventory.yaml` (using the IP from above)

1. [Deploy RPK](../QUICKSTART.md#deploy)

[Back to Quick Start](../QUICKSTART.md#dns-options)
