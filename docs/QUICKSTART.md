# Quick Start

- [Quick Start](#quick-start)
  - [TKG Installation](#tkg-installation)
  - [Cluster Sizing Requirements](#cluster-sizing-requirements)
  - [Supported Installations](#supported-installations)
  - [Administrative Access](#administrative-access)
  - [Cloud Provider Specific Steps](#cloud-provider-specific-steps)

## TKG Installation

Prior to running the commands to deploy RPK, it is required that you have a pre-existing
TKG **workload** cluster deployed.  Please follow the appropriate documentation on the VMware website at
https://docs.vmware.com/en/VMware-Tanzu-Kubernetes-Grid/1.2/vmware-tanzu-kubernetes-grid-12/GUID-tanzu-k8s-clusters-create.html in order to deploy your first TKG workload cluster.

See also [Cluster Sizing Requirements](#cluster-sizing-requirements) and [Supported Installations](#supported-installations) for sizing and installation requirements.

**NOTE:** *Deploying RPK to a management cluster is not supported at this time.*

**NOTE:** *In this current iteration, TKG* **MUST** *be deployed to an internet-connected environment. Air-gapped
deployments may be looked at to extend functionality of RPK in future releases.*

## Cluster Sizing Requirements

RPK brings an a la carte style of platform, meaning you can specify which components that you would like installed via
profiles (see [./COMPONENTS.md](COMPONENTS.md)).  Because of this reason, it is hard to determine the size of cluster that
you will need.  Each component that you select to be part of your platform has its individual requirements documented
at the component README.md file.

That being said, the platform role, which encompasses most all platform components requires the following sizing at
the time of this writing:

- 6 Worker Nodes
- 4 vCPU per Worker Node
- 16GB Memory per Worker Node
- Any Control Plane size is acceptable

## Supported Installations

RPK should technically work against most flavors of Kubernetes, however it has largely been validated against VMware's TKG
platform.  The following table outlines the common supported installations of RPK that have been tested.  Other variations
may work:

| Component    | Supported Versions |
| ------------ | ------------------ |
| Kubernetes   | 1.16 thru 1.20     |
| CNI Plugin   | Antrea or Calico   |
| TKG Versions | 1.1 thru 1.3       |

## Administrative Access

Once your TKG workload cluster is deployed, you will need to ensure that your kubeconfig file, (located at
`~/.kube/config` by default) provides administrative access to your TKG cluster.  If you setup your
first TKG workload cluster in the previous step, you will need to ensure that you have run the
`tkg get credentials <my_cluster_name>` command to ensure that the context exists for your TKG
workload cluster in your kubeconfig file.


## Cloud Provider Specific Steps

Start here by clicking your cloud provider type once you have a TKG workload cluster online.

- [AWS Instructions](providers/aws.md)
- [Azure Instructions](providers/azure.md)
- [KIND Instructions](providers/kind.md)
- [VMware Instructions](providers/vmware.md)
