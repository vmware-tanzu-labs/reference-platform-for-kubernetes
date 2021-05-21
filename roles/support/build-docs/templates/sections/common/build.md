Build a Docker image to run the Ansible playbooks.  This will relieve you of the need to install the Python dependencies locally.  To build the Docker image with the latest state of the project locally, run:

```bash
make build
```

To build this Docker image with custom image names (default: rpk) and image versions (default: latest), run the
following, being sure to substitute in your desired variables appropriately:

```bash
IMAGE=rpk-custom VERSION=v1.0.0 make build
```

> **NOTE:** *If using custom names, ensure that these match during the deploy stage.*
