#### Deploy RPK

You can deploy using the docker image with the latest state of the project locally.  To do so, run:

```bash
make deploy
```

To run when using custom image names, image versions, and/or inventory run the following, being sure to substitute your desired variables appropriately:

```bash
INVENTORY=/path/to/my/inventory.yaml IMAGE=rpk-test VERSION=v1.0.0 make deploy
```

> **NOTE:** *If you are deploying RPK to a Tanzu Mission Control provisioned cluster, you must run the following version of the command:*

```bash
TMC=true make deploy
```
