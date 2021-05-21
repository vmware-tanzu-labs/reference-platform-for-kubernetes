# Copyright 2006-2021 VMware, Inc.
# SPDX-License-Identifier: MIT
.DEFAULT_GOAL := help
.PHONY: build

# image build vars
IMAGE ?= projects.registry.vmware.com/rpk/rpk
IMAGE_VERSION ?= v1.4.0
IMAGE_BASE ?= projects.registry.vmware.com/rpk/rpk-base
IMAGE_BASE_VERSION ?= v1.4.0

# rpk vars
INVENTORY ?= `pwd`/build/inventory.yaml
KUBECONFIG_FILE ?= ~/.kube/config
CONTAINER_PREFIX ?= rpk

# provider vars
AWS_INVENTORY ?= examples/aws-inventory.yaml
VMWARE_INVENTORY ?= examples/vmware-inventory.yaml
AZURE_INVENTORY ?= examples/azure-inventory.yaml
KIND_INVENTORY ?= examples/kind-inventory.yaml
KIND_BASE_CONFIG ?= examples/kind/config.yaml
KIND_CONFIG ?= build/kind-config.yaml
KIND_CLUSTER ?= rpk

# docker vars
DOCKER_BASE_ARGS = $(CUSTOM_DOCKER_ARGS) \
	--rm -i \
	-v $(KUBECONFIG_FILE):/ansible/.kube/config:ro \
	-v $(INVENTORY):/ansible/inventory.yaml:ro \
	-v `pwd`/ansible.cfg:/ansible/ansible.cfg:ro \
	-v `pwd`/build/hosts:/etc/hosts:rw \
	--network host

DOCKER_ARGS = $(DOCKER_BASE_ARGS) \
	$(IMAGE):$(IMAGE_VERSION)

DOCKER_TEST_ARGS = $(DOCKER_BASE_ARGS) \
	-v `pwd`/roles:/ansible/roles:ro \
	-v `pwd`/lib:/ansible/lib:ro \
	-v `pwd`/build:/ansible/build:rw \
	-v `pwd`/profiles:/ansible/profiles:rw \
	-v `pwd`/site.yaml:/ansible/site.yaml:ro \
	$(IMAGE_BASE):$(IMAGE_BASE_VERSION)

RPK_ARGS = $(CUSTOM_RPK_ARGS) \
	-i "/ansible/inventory.yaml"

RPK_TEST_ARGS = $(CUSTOM_RPK_ARGS) \
	-l "-vv" \
	-i "/ansible/inventory.yaml"

# optional vars
ifdef TMC
	DOCKER_BASE_ARGS += -v ~/.vmware-cna-saas:/root/.vmware-cna-saas:ro -v ~/.vmware-cna-saas:/ansible/.vmware-cna-saas:ro
endif

ifdef AVI_OVA
	DOCKER_BASE_ARGS += -v $(AVI_OVA):/ansible/avi.ova:ro
endif

#help: @ List available tasks on this project
help:
	@echo "user tasks:"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | grep -v '{developer}' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\t\033[32m%-30s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "developer tasks:"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | grep '{developer}' | awk -F'{developer}' '{print $$1}' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\t\033[31m%-30s\033[0m %s\n", $$1, $$2}'

setup.hosts:
	@touch build/hosts

check.inventory:
	@if [[ -f $(INVENTORY) ]]; then echo "inventory file $(INVENTORY) already exists...please manually backup and re-run make setup.<PROVIDER> command"; exit 1; fi

require.role:
	@test -n "$(ROLE)"  # $$ROLE

#setup.kind: @ Create a KIND inventory file used for deploying RPK
setup.kind: check.inventory
	cp $(KIND_INVENTORY) $(INVENTORY)
	cp $(KIND_BASE_CONFIG) $(KIND_CONFIG)
	kind create cluster --config $(KIND_CONFIG) --name $(KIND_CLUSTER) --wait 10m
	@echo Please run the following before deploying RPK:
	@echo make setup.kind.network

#setup.kind.network: @ Install Calico for KIND cluster
setup.kind.network:
	kubectl apply -f examples/kind/calico.yaml --wait=true

#setup.aws: @ Create an AWS inventory file used for deploying RPK
setup.aws: check.inventory
	cp $(AWS_INVENTORY) $(INVENTORY)

#setup.vmware: @ Create a VMware inventory file used for testing or deploying RPK
setup.vmware: check.inventory
	cp $(VMWARE_INVENTORY) $(INVENTORY)

#setup.azure: @ Create an Azure inventory file used for deploying RPK
setup.azure: check.inventory
	cp $(AZURE_INVENTORY) $(INVENTORY)

#build: @ Builds local RPK image in the current state {developer}
build:
	docker build ./ -t $(IMAGE):$(IMAGE_VERSION)

#build.base: @ Builds the base RPK image in the current state {developer}
build.base:
	docker build ./ -t $(IMAGE_BASE):$(IMAGE_BASE_VERSION) -f ./Dockerfile.base

#build.docs: @ Builds RPK's documentation {developer}
build.docs:
	docker run --name $(CONTAINER_PREFIX)-build.docs \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		-v `pwd`/docs:/ansible/docs:rw \
		-v `pwd`/roles/support/build-docs/templates/sections:/ansible/roles/support/build-docs/templates/sections:rw \
		-v `pwd`/roles/support/build-docs/templates/built:/ansible/roles/support/build-docs/templates/built:rw \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS) \
		-r "build-docs" \
		-e "-e rpk_profile=single" \
		-a "support"

#build.profiles: @ Builds the base RPK component profiles after updates to components.yaml {developer}
build.profiles:
	docker run --name $(CONTAINER_PREFIX)-build.profiles \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS) \
		-r "build-profiles" \
		-a "support"

#new.role: @ Creates a new individual role (requires ROLE env variable) {developer}
new.role: require.role
	docker run --name $(CONTAINER_PREFIX)-new.role \
		-e 'RPK_ANSIBLE_ROLE=$(ROLE)' \
		--entrypoint=/ansible/scripts/make-role.sh \
		-v `pwd`/roles:/ansible/roles \
		-v `pwd`/ci/scripts:/ansible/scripts \
		$(IMAGE):$(IMAGE_VERSION)

#clean.image: @ Removes the local RPK image
clean.image:
	docker rmi $(IMAGE):$(IMAGE_VERSION)

#clean.manifests: @ Removes previously rendered templates and manifests
clean.manifests:
	rm -rf `pwd`/build/staging
	rm -rf `pwd`/build/manifests

#clean.artifacts: @ Removes the local build artifacts, including personalizations
clean.artifacts:
	rm -rf `pwd`/build/*

#clean.cluster: @ Removes all objects from the cluster
clean.cluster: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-clean.cluster \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_ARGS) \
		$(RPK_ARGS) \
		-a "clean"

#clean.test.cluster: @ Removes all objects from the cluster using current directory state {developer}
clean.test.cluster: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-clean.test.cluster \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS) \
		-a "clean"

#clean.kind: @ Removes the KIND cluster
clean.kind:
	kind delete cluster --name $(KIND_CLUSTER)

#clean.role: @ Runs cleanup tasks on a specific role using state from versioned image (requires ROLE env variable)
clean.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-clean.role \
		$(DOCKER_ARGS) \
		$(RPK_ARGS) \
		-r "$(ROLE)" \
		-a "clean"

#clean.test.role: @ Runs cleanup tasks on a specific role using current directory state (requires ROLE env variable) {developer}
clean.test.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-clean.test.role \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS) \
		-r $(ROLE) \
		-a "clean"

#deploy.role: @ Deploys an individual role and dependencies using current state from versioned image (requires ROLE env variable)
deploy.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-deploy.role \
		$(DOCKER_ARGS) \
		$(RPK_ARGS) \
		-r $(ROLE)

#deploy.test.role: @ Deploys an individual role and dependencies using current directory state (requires ROLE env variable) {developer}
deploy.test.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-deploy.test.role \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS) \
		-r $(ROLE)

#deploy: @ Deploys RPK using current state from versioned image
deploy: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-deploy \
		$(DOCKER_ARGS) \
		$(RPK_ARGS)

#deploy.test: @ Deploys RPK using current directory state {developer}
deploy.test: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-deploy.test \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_ARGS)

#demo.role: @ Runs demo tasks on a specific role using current state from versioned image (requires ROLE env variable)
demo.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-demo.role \
		$(DOCKER_ARGS) \
		$(RPK_ARGS) \
		-r "$(ROLE)" \
		-a "demo"

#demo.test.role: @ Runs demo tasks on a specific role using current directory state (requires ROLE env variable) {developer}
demo.test.role: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-demo.test.role \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS)    \
		-r "$(ROLE)" \
		-a "demo"

#infra: @ Runs infrastructure setup tasks on a specific role using current state from versioned image (requires ROLE env variable)
infra: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-infra \
		$(DOCKER_ARGS)          \
		$(RPK_ARGS)             \
		-r "$(ROLE)/infra"

#infra.test: @ Runs infrastructure setup tasks on a specific role using current directory state (requires ROLE env variable) {developer}
infra.test: require.role setup.hosts
	docker run --name $(CONTAINER_PREFIX)-infra.test \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		$(DOCKER_TEST_ARGS) \
		$(RPK_TEST_ARGS)    \
		-r "$(ROLE)/infra"

#lint.ansible: @ Runs linting functions against Ansible {developer}
lint.ansible: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-lint.ansible \
		--entrypoint='/ansible/scripts/lint-ansible.sh' \
		-v `pwd`/site.yaml:/ansible/site.yaml:ro \
		-v `pwd`/ci/scripts:/ansible/scripts \
		-v `pwd`/.ansible-lint-rules:/ansible/.ansible-lint-rules:ro \
		-v `pwd`/.ansible-lint:/ansible/.ansible-lint:ro \
		$(DOCKER_TEST_ARGS)

#lint.dirs: @ Runs linting functions to check directories {developer}
lint.dirs: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-lint.dirs \
		--entrypoint='/ansible/scripts/lint-directories.sh' \
		-v `pwd`/ci/scripts:/ansible/scripts \
		$(DOCKER_TEST_ARGS)

#lint.files: @ Runs linting functions to check file extensions {developer}
lint.files: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-lint.files \
		--entrypoint='/ansible/scripts/lint-file-extensions.sh' \
		-v `pwd`/ci/scripts:/ansible/scripts \
		$(DOCKER_TEST_ARGS)

#lint.yaml: @ Runs linting functions against YAML files {developer}
lint.yaml: setup.hosts
	docker run --name $(CONTAINER_PREFIX)-lint.yaml \
		--entrypoint='/ansible/scripts/lint-yaml.sh' \
		-v `pwd`/.yamllint:/ansible/.yamllint:ro \
		-v `pwd`/ci/scripts:/ansible/scripts \
		$(DOCKER_TEST_ARGS)

#lint.all: @ Runs all lint.* targets {developer}
lint.all: lint.dirs lint.files lint.yaml lint.ansible

#list.profiles: @ Lists all available RPK profiles
list.profiles:
	@docker run \
	    --rm \
		--entrypoint=/bin/bash \
		$(IMAGE):$(IMAGE_VERSION) \
		-c 'printf "All available RPK profiles:\n\n`ls -1 /ansible/profiles/ | cut -d "." -f 1 | egrep -v "components|single" | sort -u`\n"'

#list.components.all: @ Lists all available RPK components
list.components.all:
	@docker run \
	    --rm \
		--entrypoint=/bin/bash \
		$(IMAGE):$(IMAGE_VERSION) \
		-c 'printf "All available RPK Components:\n\n" && yq r /ansible/profiles/components.yaml "rpk_components(*).name" | sort -u'

#list.components.advanced: @ Lists all available RPK components from the "advanced" profile
list.components.advanced:
	@docker run \
	    --rm \
		--entrypoint=/bin/bash \
		$(IMAGE):$(IMAGE_VERSION) \
		-c "printf \"All available RPK Components in 'advanced' profile:\n\n\" && yq r /ansible/profiles/advanced.yaml 'rpk_components(*).name' | sort -u"

#list.platform: @ Lists all available RPK components from the "platform" profile
list.components.platform:
	@docker run \
	    --rm \
		--entrypoint=/bin/bash \
		$(IMAGE):$(IMAGE_VERSION) \
		-c "printf \"All available RPK Components in 'platform' profile:\n\n\" && yq r /ansible/profiles/platform.yaml 'rpk_components(*).name' | sort -u"

# TODO: validate targets incomplete at the current moment
#validate.role: @ Validates an individual role using sonobuoy plugin (requires ROLE env variable)
validate.role: require.role
	sonobuoy run -p $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml

#validate.status: @ Gets status of current validation test
validate.status:
	sonobuoy status

#validate.clean: @ Deletes sonobuoy assets used for validation test
validate.clean:
	sonobuoy delete

#validate.retrieve @ Retrieves test results and writes them to /tmp directory
validate.retrieve:
	rm -rf /tmp/rpk-validation-results
	mkdir /tmp/rpk-validation-results
	sonobuoy retrieve /tmp/rpk-validation-results

#validate.results @ Displays results summary
validate.results:
	sonobuoy results /tmp/rpk-validation-results/2*

#validate.detail @ Displays contents of results file generated by sonobuoy plugin (requires ROLE env variable)
validate.detail: require.role
	tar xf /tmp/rpk-validation-results/2* -C /tmp/rpk-validation-results/
	cat /tmp/rpk-validation-results/plugins/rpk-$(ROLE)/sonobuoy_results.yaml

#validate.logs @ Displays contents of plugin logs
validate.logs: require.role
	cat /tmp/rpk-validation-results/plugins/rpk-$(ROLE)/results/global/plugin.log

#validate.test.role @ Tests validation for role using sonobuoy plugin (requires ROLE and IMAGE env variables)
validate.test.role: require.role
	@test -n "$(VALIDATE_TEST_IMAGE)"
	docker build -t $(VALIDATE_TEST_IMAGE) $(PWD)/roles/$(ROLE)/validate/
	docker push $(VALIDATE_TEST_IMAGE)
	awk -v NEWIMAGE="  image: $(VALIDATE_TEST_IMAGE)" '{if ( $1 /^ *image: / ) print NEWIMAGE; else { print $0 }}' $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml > /tmp/sonobuoy-plugin.yaml
	mv /tmp/sonobuoy-plugin.yaml $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml
	sonobuoy run -p $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml

#validate.release.role @ Cut a new release of validation testing container image for a role (requires ROLE and VALIDATE_VERSION env variables)
validate.release.role: require.role
	@test -n "$(VALIDATE_VERSION)"
	docker build -t cloudnativereadiness/rpk-$(ROLE)-validate:$(VALIDATE_VERSION) $(PWD)/roles/$(ROLE)/validate/
	docker push cloudnativereadiness/rpk-$(ROLE)-validate:$(VALIDATE_VERSION)
	awk -v NEWIMAGE="  image: cloudnativereadiness/rpk-$(ROLE)-validate:$(VALIDATE_VERSION)" '{if ( $1 /^ *image: / ) print NEWIMAGE; else { print $0 }}' $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml > /tmp/sonobuoy-plugin.yaml && mv /tmp/sonobuoy-plugin.yaml $(PWD)/roles/$(ROLE)/validate/sonobuoy-plugin.yaml
