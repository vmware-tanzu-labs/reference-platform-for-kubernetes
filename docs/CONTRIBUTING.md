# Contributing to RPK

The following is a set of guidelines for contributing to RPK.

## Table of Contents

- [Contributing to RPK](#contributing-to-rpk)
  - [Table of Contents](#table-of-contents)
  - [General](#general)
    - [Relation to Tanzu Developer Center](#relation-to-tanzu-developer-center)
  - [Git](#git)
    - [Developer Workflow](#developer-workflow)
    - [Branches > Forks](#branches--forks)
    - [Development Environment](#development-environment)
      - [Development Environment Prerequisites](#development-environment-prerequisites)
      - [Developing with a Docker Image](#developing-with-a-docker-image)
      - [Developing with a Local Environment](#developing-with-a-local-environment)
  - [Ansible](#ansible)
    - [Why Ansible?](#why-ansible)
    - [TDC Module == Ansible Role/Component](#tdc-module--ansible-rolecomponent)
    - [Component Structures](#component-structures)
    - [Profiles](#profiles)
      - [RPK's Default Profiles](#rpks-default-profiles)
        - [Listing Available Profiles](#listing-available-profiles)
        - [Listing All Available RPK Components](#listing-all-available-rpk-components)
        - [Listing Components in the Platform Profile](#listing-components-in-the-platform-profile)
        - [Listing Components in the Advanced Profile](#listing-components-in-the-advanced-profile)
      - [Profile Definitions](#profile-definitions)
        - [Example Component Specification](#example-component-specification)
        - [Profile Component Dictionary Field Reference](#profile-component-dictionary-field-reference)
        - [Building and Rebuilding Profiles](#building-and-rebuilding-profiles)
    - [Editing Cloud Provider QUICKSTART Documentation](#editing-cloud-provider-quickstart-documentation)
      - [Build RPK Documentation](#build-rpk-documentation)
    - [New Roles](#new-roles)
      - [Assign your Component to a Profile](#assign-your-component-to-a-profile)
      - [Rebuild RPK Profiles](#rebuild-rpk-profiles)
      - [Rebuild RPK Documentation](#rebuild-rpk-documentation)
    - [Component Dependencies](#component-dependencies)
    - [Prefer Simplicity Over Complexity](#prefer-simplicity-over-complexity)
  - [Questions?](#questions)

## General

### Relation to Tanzu Developer Center

The goal of RPK is to provide a reference implementation platform which encompasses the
ideas and best practices as captured in the Tanzu Developer Center Kubernetes Guides (https://tanzu.vmware.com/developer/guides/kubernetes/).  Any
contributions to RPK should align with the ideals and practices that are found there.

For more information about RPK and it's relation to TDC, view the following links:

* [Overview](./OVERVIEW.md)
* [Architecture](./ARCHITECTURE.md)

## Git

### Developer Workflow

The following outlines the simple developer workflow that we use for RPK.

IMAGE HERE

1. Contributor discovers a bug or has an idea for a new feature or an improvement to the existing processes.
2. Contributor opens an issue and fills out all of the appropriate issue from the template.
3. Issue is either approved, denied, or more information is requested.
   1. Upon denial, user may attend the [Standup Meetings](#standup-meetings) or chat in the [Slack Communications Channel](#slack-communications-channel) to get issue approved.
   2. Upon request for more information, issue will remain pending until approved.  Information should be added to the issue directly.
   3. **NOTE:** failure to follow this process will result in your ideas, code, etc. to be revoked from admission into the repository.
4. Contributor applies the `In Progress` label to the issue and assigns himself/herself to the issue.
5. Contributor clones the RPK repo (e.g. `git clone git@github.com/vmware-tanzu-labs/reference-platform-kubernetes.git`).  See [Development Environment](#development-environment) for information on environment setup.
6. Contributor creates a new branch (`e.g. git checkout -b my-cool-new-feature`) on their local development workstation.
7. Working in the new branch on the local development workstation, the contributor modifies the code needed to address the opened and approved issue.
8. Contributor commits and pushes the changes to the RPK repo (`e.g. git add*; git commit -a -m 'Fixes #1, my commit message'; git push --set-upstream origin my-cool-new-feature`)
9. Contributor opens a merge request in GitLab and fills out the appropriate information in the Merge Request.
10.  A CI pipeline is kicked off.  See [PIPELINE.md](PIPELINE.md) for more details.
   1. **NOTE:** failed CI pipeline runs will not be merged.
   2. **NOTE:** please keep commits to their individual modules (e.g. `container-registry`, or `storage`) as this helps unit test the independent modules.
11. Once CI pipeline passes, you may remove the `WIP` from your merge request.
12. If additional changes are requested, steps 7-8 can be repeated until the branch is approved for merge by the maintainers.
13. Once the request approved, your code is merged!

### Branches > Forks

For now, we prefer branches over forks.  This allows us to do a couple things:

1. **CI Pipeline Testing:**  This doesn't work with forks.
2. **Maintainer Edits:**  Sometimes we want to make simple, quick code changes.

**NOTE:** if you open a merge request from a fork, you will be asked to move to a branch instead!

### Development Environment

#### Development Environment Prerequisites

1. [Install TKG](QUICKSTART.md#tkg-installation)
2. [Provide admin access to TKG](QUICKSTART.md#administrative-access)
3. [Choose a DNS Solution](QUICKSTART.md#dns-options)
4. [Choose a Cloud Provider](QUICKSTART.md#provider-specific-steps)
5. Set inventory variables based on your Cloud Provider:
  * [AWS](providers/aws.md#set-variables)
  * [Azure](providers/azure.md#set-variables)
  * [KIND](providers/kind.md#set-variables)
  * [VMware](providers/vmware.md#set-variables)

#### Developing with a Docker Image

Ansible requires several python dependencies in order to run RPK.  Often times, this
presents challenges when porting between different operating systems and python versions.  To
get around these issues, you can build RPK in a Docker image.  To do so, run:

```bash
make build
```

For custom names and versions, run:

```bash
IMAGE=rpk-custom VERSION=v1.0.0 make build
```

To test an individual role using your previously built Docker image:

```bash
ROLE=my-role make deploy.test.role
```

The above is run in a CI pipeline each time you push your code for only the modules that are updated.  For this
reason, **please keep your commits and pushes independent to each module**!!!

To test a full end-to-end deployment using your previously built Docker image, also using the current state
of the project, run:

```bash
make deploy.test
```

#### Developing with a Local Environment

> :warning: this process is not formally supported, but may be wanted for developers who want to use
`ansible-playbook` directly.

An alternative to developing with the Docker image is to install dependencies to your development
machine directly.  This avoids having to rebuild the image each time your `requirements.txt` file
changes, but could make it difficult to develop across different platforms.

To setup your local development environment, first set up venv.  This will allow you to install
python packages for RPK into a virtual environment, rather than system-wide:

```bash
python3 -m venv ansible-virtualenv
```

Activate the venv:

```bash
. ansible-virtualenv/bin/activate
```

Install python dependendencies:

```bash
pip3 install -r requirements.txt
```

Run the playbook against a single role:

```bash
./bin/rpk -r MY_ROLE_NAME
```

> The above roughly translates into: `ansible-playbook test.yaml -e 'rpk_role_name=MY_ROLE_NAME' -c local -i build/inventory.yaml`

Run a full deployment of RPK:

```bash
./bin/rpk
```

The above roughly translates into: `ansible-playbook site.yaml -c local -i build/inventory.yaml --skip-tags module_dependencies`


## Ansible

### Why Ansible?

A common question that gets asked is "Why the use of declarative automation tool to apply declarative manifests in Kubernetes?".  The
answer is simple.  We wanted to quickly give a reference implementation of a platform using best practices from TDC and we
wanted to quickly "make it exist".  Using the talent that we had, Ansible was the easiest, shortest path to success.

Long-term plans for Ansible are still up in the air.

### TDC Module == Ansible Role/Component

In RPK, we now refer to Ansible roles as components.

For each module from the Tanzu Developer Center Kubernetes Guides (https://tanzu.vmware.com/developer/guides/kubernetes/), we tie that into an Ansible [role](https://docs.ansible.com/ansible/latest/user_guide/playbooks_reuse_roles.html) (RPK Component).  This
allows us the ability to layer on the components as building blocks with our existing processes.

### Component Structures

For each component, we require the following structure:

```
├── clean
│   └── tasks
│       └── main.yaml
├── common
│   └── defaults
│       └── main.yaml
├── defaults
│   └── main.yaml
├── demo
│   ├── tasks
│   │   └── main.yaml
├── tasks
│   └── main.yaml
├── templates
├── README.md
├── .dependencies.yaml
```

The following are deviances from common Ansible best practices:

1. **Use of common/vars/ for variables**.  This allows us to require the variables in other roles without having to statically define variables.
2. **Use of .dependencies.yaml for module dependencies**.  This structure defines *ONLY* the modules/roles in which each role is dependent upon.
3. **Use of clean/ sub-role**.  This role is the code to cleanup the Ansible role using `ROLE=my-role make clean.role`.
4. **Use of demo/ sub-role**.  This role is the code to demonstrate the role (either print out info or create K8S objects) using `ROLE=my-role make demo.role`.

### Profiles

Profiles were added to RPK to address differing needs:
  1. Support to deploy all of Tanzu Advanced features in a single pass
  2. A means for grouping RPK components that should belong together
  3. To support multi-cluster deployments where an end-user may want to provide differing platform services in differing clusters

#### RPK's Default Profiles

RPK provides two base profiles:
  1. `platform` => this is the default
  2. `advanced` => this is used to install Tanzu Advanced components

##### Listing Available Profiles

To view available profiles, run:

```bash
make list.profiles
```

##### Listing All Available RPK Components

To view all available RPK components, run:

```bash
make list.components.all
```

##### Listing Components in the Platform Profile

To view available components in the `platform` profile in the order they will be deployed, run:

```bash
make list.components.platform
```

##### Listing Components in the Advanced Profile

To view available components in the `advanced` profile in the order they will be deployed, run:

```bash
make list.components.advanced
```

#### Profile Definitions

Today, `components.yaml` can be found in the `profiles` [directory](../profiles).  This file defines several things about each RPK component:
  1. **Order of deployment.**  The order of each RPK component in the list is exactly how it would execute during that profile's deployment.

  1. **It's name**, which must be the directory name within the project's `roles` directory.

  1. Whether or not it is **enabled** by default.  We provide some components that can be deployed out of band as desired, but set them to disabled for a normal deployment.

  1. Whether or not it has a **demo** component.

  1. A list of **profiles** that the component applies to.  A component can belong to as many profiles as you would like.  We currently do not support custom profiles, but we are looking to change this in a future release.

  **NOTE:** *A component can have dependencies that don't typically run in a profile.  This is fine so long as all of a component's dependencies have been listed in the component's `.dependencies.yaml`. See [Component Dependencies](#component-dependencies) for more details.*

##### Example Component Specification

```yaml
- name:    "workload-tenancy"
  enabled: true
  profiles:
    - "platform"
```

##### Profile Component Dictionary Field Reference

These values can be set in `profiles/components.yaml`

| Variable Name | Description | Required | Type |
| --- | --- | --- | --- |
| name | Name of the RPK component (Ansible Role) directory relative to the `roles` directory. Example: `roles/workload-tenancy`.  The name would be `workload-tenancy`. | yes | string |
| enabled | Whether the component is enabled or not | yes | boolean |
| profiles | List of profiles the component applies to | yes | list |

##### Building and Rebuilding Profiles

Once a component has been modified or created in `components.yaml`, the profiles must be rebuilt.  `components.yaml` should be the only place you make changes, as it is the specification for how all of the RPK profiles should be built.

To rebuild the RPK profiles from an updated `components.yaml`, run `make build.profiles`.  If you are developing a new role, you will want to add the changed profile files to your branch, commit, and push the changes.

### Editing Cloud Provider QUICKSTART Documentation

The documentation found in `/docs/providers` has been moved into a templated format.  Do not edit the files in `/docs/providers` directly, or you will find your changes get overwritten eventually.

To properly edit these files, you can edit the common documentation sections in [/roles/support/build-docs/templates/sections/common](../roles/support/build-docs/templates/sections/common).

To properly edit the cloud provider documentation templates, you can edit the files in [/roles/support/build-docs/templates/sections/main](../roles/support/build-docs/templates/sections/main).

To properly edit the title's of the cloud provider documentation, you can edit the template files (by provider name) in [/roles/support/build-docs/templates](../roles/support/build-docs/templates).

Further details can be found in [support/build-docs README.md](../roles/support/build-docs/README.md)

#### Build RPK Documentation

Once you have modified any of the template files or common sections for the cloud provider documentation, you must rebuild the templated documentation.

This is performed by running the following command:

```bash
make build.docs
```

Add and commit the changed files into the repository.

### New Roles

Whenever a new role is required, use the command `ROLE=my-new-role make new.role`, replacing `my-new-role` with your desired RPK component name.

This will provide a new component structure consistent with the remainder of the project.

#### Assign your Component to a Profile

See [Profile Definitions](#profile-definitions)

#### Rebuild RPK Profiles

This should be completed once you have assigned your component to a profile in `components.yaml`.

See [Building and Rebuilding RPK Profiles](#building-and-rebuilding-profiles)

#### Rebuild RPK Documentation

This is necessary, because the process rebuilds the file and tables at [roles/support/build-docs/templates/sections/additional_vars.md](../roles/support/build-docs/templates/sections/common/additional_vars.md) based off of the profile variables, which you have just rebuilt after assigning the new role into `components.yaml`.

See [Build RPK Documentation](#build-rpk-documentation).


### Component Dependencies

The use of the standard Ansible Role `meta/main.yaml` file to resolve component dependencies has been deprecated in RPK.  Instead, dependencies are defined via a
`.dependencies.yaml` file within each role.  Dependencies are resolved with the following process:

1. When `rpk_profile` is set to `single` (set when requesting an individual role), dependency resolution is turned on.

1. When `rpk_profile` is set to anything but `single`, dependencies are not resolved and the RPK deployer is trusting that all
dependencies are pre-defined in the profile requested in their proper order.

1. Once RPK has detected that the profile is `single`, it passes the component into a custom Ansible filter `rpk_resolve_deps`.

1. The custom filter will resolve the dependencies based on the `.dependencies.yaml` file within the requested role to ensure all required components have been deployed within the cluster first.

1. The dependencies are stored in the `rpk_components` variable, as resolved in the previous step.

1. The `site.yaml` file loops through `rpk_components` and runs each individual component.

This method has a couple limitations:

The use of `meta/main.yaml` is to resolve dependencies has been deprecated in RPK.  Instead, dependencies are defined via a
`.dependencies.yaml` file within each role.  Dependencies are resolved with the following process:

1. When `rpk_profile` is set to `single` (set when requesting an indvididual role), dependency resolution is turned on.
2. When `rpk_profile` is set to anything but `single`, dependencies are not resolved and the RPK deployer is trusting that all
dependencies are pre-defined in the profile requested in their proper order.
3. Once RPK has detected that the profile is `single`, it passes the component into a custom Ansible filter `rpk_resolve_deps`.
4. The custom filter will resolve the dependencies based on the `.dependencies.yaml` file within the requested role.
5. The dependencies are stored in the `rpk_components` variable, as resolved in the previous step.
6. The `site.yaml` file loops through `rpk_components` and runs each individual component.

This method has a couple limitations:

1. It is not a native Ansible-supported methodology like `meta/main.yaml` is.  We went with this methodology because the out-of-the box
method would re-run dependencies when requesting a single role, which significantly increased the time it took to deploy a role.
2. Child dependencies are not resolved.  This means that you must specify ALL dependencies you need for a given role and not rely on the
child role to resolve a dependency for you.  This methodology is preferred because it leaves a documentation trail of what components get
installed with each role and greatly simplifies the logic needed to compile the roles needed.

### Prefer Simplicity Over Complexity

Given the long-term roadmap for the project's use of Ansible as a deployment tool is questionable until we evaluate further, please prefer simplicity over complexity.  There are many
cool and crazy things that we can do with Ansible, however we prefer to simply use Ansible to do the following:

1. Apply Kubernetes manifests using the `common/manifest` role.
2. Talk to APIs, generally during role demonstration using `ROLE=my-role make demo.role`.
   1. Sometimes we also talk to APIs during role deployment, but this should be limited.

## Questions?

Please reach out in the [Standup Meetings](#standup-meetings) or chat in the [Slack Communications Channel](#slack-communications-channel)!!!
