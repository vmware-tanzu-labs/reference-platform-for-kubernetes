### 1. Create Cluster Inventory

A custom inventory will need to be created in order to define which cluster(s) you will be
deploying RPK to.  This inventory file is a [standard Ansible inventory file, in YAML
format](https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#inventory-basics-formats-hosts-and-groups).  The inventory file additionally contains all of the necessary variables that feed into the deployment process.

{% if provider == "kind" %}
To deploy the KIND cluster (tested with `v0.9.0`), inventory file, and CNI plugin, run the following commands:
{% else %}
To create your inventory file, run the following command:
{% endif %}

```bash
make setup.{{ provider }}
{% if provider == "kind" %}
make setup.kind.network
{% endif %}
```

> **NOTE:** The command above will drop your inventory file example at `build/inventory.yaml`.  If you need the inventory file at another location,
> you can run `INVENTORY=/path/to/my/inventory.yaml make setup.aws` for example.  Be sure that the `INVENTORY` variable matches this for following
> build and deploy steps.
