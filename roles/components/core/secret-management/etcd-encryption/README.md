# secret_management_etcd_encryption

An Ansible Role to deploy ETCD encryption at rest.  This Role assumes that the inventory target will be all control plane nodes running both ETCD and kube-apiserver containers.

If an encryption key is already in place, then the new encryption key will be prepended on top of the existing provider list. See https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/#understanding-the-encryption-at-rest-configuration for details of how this works.

## Architecture

RPK's architecture can be found [here](../../../docs/ARCHITECTURE.md#etcd-encryption).


## Default Variables

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| workload_cluster_namespace | The management cluster's namespace where the target cluster belongs to | "default" | string | yes |
| etcd_encryption_conf_path | Full path to deploy the ETCD encryption configuration. | "{{ etcd_encryption_conf_base_path }}/encryption-config.yaml" | string | yes |
| k8s_apiserver_manifest_path | Path to the kube-apiserver manifest configuration.  This assumes a `kubeadm` style manifest | "/etc/kubernetes/manifests/kube-apiserver.yaml" | string | yes |
| etcd_cert_dir | Path to the ETCD certificates directory, and is leveraged for ETCD Backups during deployment and removal of the ETCD encryption configuration | "/etc/kubernetes/pki/etcd" | string | yes |


## Example secret_management_etcd_encryption Playbook Usage

This is the expected normal usage of the Role:

```yaml
---
- name: "Configuration of ETCD encryption at rest"
  hosts: "control_plane"
  gather_facts: true
  tasks:
    - name: "Deploy ETCD encryption at rest"
      include_role:
        name: "secret_management_etcd_encryption"
```

The above example also works in the case of key rotation. The new key is generated and prepended to the list of secret keys. Re-encrypting all secrets is out of scope for this role, but can be done in a single command as follows:

```bash
kubectl get secrets --all-namespaces -o json | kubectl replace -f -
```

## Test
You can test this role by running the [test/test.sh](test/test.sh) script. The [test/inventory](test/inventory) file should be modified to suit your environment prior to running the test.

## Uninstall
If you wish to revert the changes performed by this role, there is an uninstall script [test/uninstall.sh](test/uninstall.sh).  You need to modify the [test/inventory](test/inventory) to suit your environment prior to running the uninstall.

The uninstall procedure will perform a backup of the ETCD encryption configuration, as well as the ETCD database prior to making modifications.  ***Only perform this if absolutely required as it could cause irreparable harm to your ETCD data during the decryption process.  You've been warned :)***

## Author(s)
[Andrew J. Huffman](mailto:ahuffman@vmware.com)
[George Goh](mailto:gohge@vmware.com)
