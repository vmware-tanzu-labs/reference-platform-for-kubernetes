# Setting DNS Resolvers

## MacOS
### Via UI

See: [Specify a DNS server on Mac](https://support.apple.com/en-is/guide/mac-help/mchlp2720/mac)

### Using /etc/resolver Configurations

1. Replacing `yourdomain.com` with the domain you set in `tanzu_ingress_domain`:
```bash
cd /etc/resolver
sudo touch yourdomain.com
```

2. Edit the file replacing `yourdomain.com` with the domain you set as the `tanzu_ingress_domain`, and the `1.2.3.4` IP Address with that of your DNS server: `sudo vi /etc/resolver/yourdomain.com`:
```bash
search yourdomain.com
nameserver 1.2.3.4
```


## Linux
1. Edit `/etc/resolv.conf`:
```bash
sudo vi /etc/resolv.conf
```

2. Replacing `yourdomain.com` with the domain you set as the `tanzu_ingress_domain`, and the `1.2.3.4` IP Address with that of your DNS server:
```bash
search yourdomain.com
nameserver 1.2.3.4
```
