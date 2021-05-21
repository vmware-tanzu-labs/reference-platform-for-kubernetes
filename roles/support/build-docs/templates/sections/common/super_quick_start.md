## Super Quick Start

1. `make setup.{{ provider }}` - create the inventory file with global variables.
{% if provider == "kind" %}

1. `make setup.kind.network` - deploy Calico as a CNI for KIND.

1. `kind get kubeconfig --name rpk >> ~/.kube/config` - ensure the KIND kubeconfig and contexts exist in the kubeconfig file.
{% endif %}

1. `vim build/inventory.yaml` - edit and update the variables.

1. `make build` - build an image for deploying RPK.

1. `make deploy` - deploy RPK.
