If your cluster is a TMC-provisioned cluster, the following steps need to be performed:

1. Install the TMC binary.  See https://docs.vmware.com/en/VMware-Tanzu-Mission-Control/services/tanzumc-using/GUID-7EEBDAEF-7868-49EC-8069-D278FD100FD9.html.

1. `tmc login`

1. When launching the `make deploy` command in the steps following, ensure that you run it with `TMC=true`

For example:

```bash
TMC=true make deploy
```
