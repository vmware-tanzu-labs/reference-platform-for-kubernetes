# Networking

Because TKG deploys the Container Network Interface (CNI) as part of the deployment, this
role largely will simply check to see which CNI is installed, tweak some knobs, and
secure the network based upon the CNI which is provisioned for us.


## Architecture

RPK's architecture can be found [here](../../docs/ARCHITECTURE.md#container-networking).


## Resource Sizing Requirements

The following sizing requirements must be met for this role to operate properly.  Sizing includes additional [dependencies](#dependencies).

| vCPU | Memory | Storage |
| --- | --- | --- |
| N/A | N/A | N/A - no persistent storage required |

**NOTE:** no resources required.  CNI components are deployed by TKG.


## Variables

### Default Variables

The following variables are defaulted in the `common/vars/main.yaml` file and require no additional user input.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_networking.namespace | Namespace for tanzu_networking components | "tanzu-networking" | string | yes |
| tanzu_networking.staging_dir | Local directory to write the staging manfiests to | "{{ rpk_staging_dir }}/tanzu-networking" | string | yes |
| tanzu_networking.ipip_mode | The IP-in-IP mode to use for the CNI (applicable to Calico) | "CrossSubnet" | string | yes |
| tanzu_networking.antrea_tunnel_protocol | Protocol used for network overlay tunnel (applicable to Antrea) | "geneve" | string | yes |
| tanzu_networking.antrea_crd_policy | CRD Policy type used for network policy engine (applicable to Antrea) | "ClusterNetworkPolicy" | string | yes |

### Additional Variables

The following variables must be set for proper operatiion of the role.  Variables are generally set in the variables file
at `build/inventory.yaml` of the root of this project.

| Variable Name | Description | Default Value | Variable Type | Required |
| --- | --- | --- | --- | --- |
| tanzu_kubectl_context | Name of context to use for connection to Kubernetes | - | string | yes |


## Dependencies

Also see `.dependencies.yaml` to view role dependencies which are run when running the role
independently.

* NONE


## Deploying

**NOTE:** roles from `.dependencies.yaml` are also deployed.

In order to deploy the role from a versioned image:

```bash
ROLE=networking make deploy.role
```

If you've made changes to the role and need to verify your changes:

```bash
ROLE=networking make deploy.test.role
```


## Demonstrating

Once the role has run successfully, you should be able to demonstrate the role.

In order to demonstrate the role:

**NOTE:** this is a work in progress.  See [Verification](#verification) for steps to manually demonstrate
the role.

```bash
ROLE=networking make demo.role
```


### Verification

This role should have created a global deny-all network policy in your cluster, excluding any TMC namespace, kube-system and tanzu-*.

1) Change the context to the cluster to which you wish to check:

`kubectl config use-context <my_context_name>`

2) View the Global Network Policy:

`kubectl get globalnetworkpolicy`

To check and make sure the Global Network Policy is working:

1) Change the context to the cluster to which you wish to check:

`kubectl config use-context <my_context_name>`

2) Create a test deployment:

`kubectl apply -f ./test/deployment.yaml`

3) Expose the deployment (some mileage may vary depending upon your setup):

`kubectl expose deployment nginx --type=LoadBalancer --name=nginx-service`

4) Find the IP or hostname of the exposed service:

`kubectl describe service nginx-service`

5) Check to ensure the port is in accessible:

curl http://172.26.0.202
curl: (7) Failed to connect to 172.26.0.202 port 80: Connection refused

6) If you delete the global network policy, you will again get a successful connection to the service.

## Cleaning

To remove the role, from the root of the repo:

```bash
ROLE=networking make clean.role
```


## Author(s)
[Dustin Scott](mailto:sdustin@vmware.com)
