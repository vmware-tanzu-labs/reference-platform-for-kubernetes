## Ingress TLS Certificate Options

RPK uses [cert-manager](https://cert-manager.io/) to manage TLS certificates for Ingress resources.

### Default

By default a new Certificate Authority is created and all Ingress certificates are signed by that CA and you do not need to change anything.

After deploying RPK the best way to download the CA Cert is by logging into Harbor and clicking through to the `library` project. You will find a button to download the CA Cert which you can then use with Docker, or add to your Browser. You can also find it in a secret called `ca-keypair` in the `tanzu-security` namespace.

The Caveat to using the defaults is that Kubernetes needs to trust the created root CA. To do this we have Kubernetes Jobs that will copy the CA cert to the Kubernetes nodes and reconfigure them to be trusted. This is good enough for testing or a POC, but is not robust enough for production.

We suggest you either provide a pre-created Wildcard Certificate Pair, your own CA Pair, or use [Let's Encrypt](https://letsencrypt.org/).

### Let's Encrypt

If you have public Ingress IPs and DNS you can modify your Inventory to get cert-manager to  negotiate certificates with [Let's Encrypt](https://letsencrypt.org/) and your certificates will be automatically trusted.

```yaml
tanzu_security:
  tls_providers: ["ca", "letsencrypt-stage", "letsencrypt-prod"]
  default_tls_provider: "letsencrypt-prod"
```

### Provide your own Wildcard Cert

You can choose to provide RPK with a pre-created Wildcard Certificate to use for Ingress.
You will be responsible for adding the CA certificate to your Kubernetes nodes (see [TKG Docs](https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-tanzu-k8s-clusters-config-plans.html?hWord=N4IghgNiBcIMYFcDOAXA9gWwARwKYCcUBLAMyLjBVyzARQAs18iUBPEAXyA#custom-certificate-authority-vsphere-8) for an example of doing so with TKG). If you do so you will want to disable updating the Kubernetes Nodes with the CA, and provide the ca cert, cert, and key for the wildcard:

```yaml
tanzu_security:
  actions:
    update_k8s_ca: false
  tls_providers: ["wildcard", "ca"]
  default_tls_provider: "wildcard"
  wildcard:
    tls_root_ca_cert: <base64 encoded CA cert>
    tls_cert: <base64 encoded wildcard cert>
    tls_key: <base64 encoded wildcard cert>
```

### Provide your own CA Cert

You can choose to provide RPK with a pre-created CA cert/key pair to use to create TLS certificates as needed for Ingress. You will be responsible for adding the CA certificate to your Kubernetes nodes and provide the ca cert and ca key:

Create a CA cert (feel free to use your preferred certificate tool of choice):

```bash
docker run -v /tmp/certs:/certs -e SSL_SUBJECT=rpk-ca-cert paulczar/omgwtfssl
```

Append the following to the TKG ytt overlay for your provider (ex `~/.tkg/providers/infrastructure-aws/ytt/aws-overlay.yaml`) replacing the Certificate in both places with the contents of `/tmp/certs/ca.pem`:

> Note: This must be done before creating your cluster.

```yaml
#! Add and trust your custom CA certificate on all Control Plane nodes.
#@overlay/match by=overlay.subset({"kind":"KubeadmControlPlane"})
---
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: KubeadmControlPlane
spec:
  kubeadmConfigSpec:
    #@overlay/match missing_ok=True
    preKubeadmCommands:
    #@overlay/append
    - |
      cat <<EOF > /etc/ssl/certs/myca.pem
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
      EOF
      #@overlay/append
    - openssl x509 -in /etc/ssl/certs/myca.pem -text >> /etc/pki/tls/certs/ca-bundle.crt
    #@overlay/append
    - c_rehash

#! Add and trust your custom CA certificate on all worker nodes.
#@overlay/match by=overlay.subset({"kind":"KubeadmConfigTemplate"})
---
spec:
  template:
    spec:
      #@overlay/match missing_ok=True
      preKubeadmCommands:
      #@overlay/append
      - |
        cat <<EOF > /etc/ssl/certs/myca.pem
        -----BEGIN CERTIFICATE-----
        ...
        -----END CERTIFICATE-----
        EOF
      #@overlay/append
      - openssl x509 -in /etc/ssl/certs/myca.pem -text >> /etc/pki/tls/certs/ca-bundle.crt
      #@overlay/append
      - c_rehash
```

Create your Kubernetes cluster:

```bash
tkg create cluster tkg-20210106 --plan=dev --worker-machine-count 3 --worker-size m5.2xlarge
```

Update your inventory to include the following variables:

```yaml
tanzu_security:
  actions:
    update_k8s_ca: false
  ca:
    # cat /tmp/certs/ca.pem | base64
    tls_root_ca_cert: <base64 encoded CA cert>
    # cat /tmp/certs/ca-key.pem | base64
    tls_root_ca_key: <base64 encoded CA key>
```

Deploy RPK as usual.
