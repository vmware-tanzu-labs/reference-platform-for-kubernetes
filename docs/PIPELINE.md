# CI/CD Pipeline

RPK utilizes a CI/CD Pipeline in order to ensure the quality of the code being submitted.
The pipeline runs inside several different Gitlab Runners in order to perform different task to
test the code.

## Test Cases

Although not all test cases are possible, due to the near infinite different combinations, we try and
provide a significant amount of test cases that we will encounter in the field.

The following test scenarios are covered by the CI/CD Pipeline:

| Test Case          | AWS     | VMware   | KIND   | Azure | Comments                                                                         |
| ------------------ | ------- | -------- | ------ | ----- | -------------------------------------------------------------------------------- |
| Kubernetes Version | 1.18.n  | 1.19.n   | N/A    | N/A   |                                                                                  |
| TKG Version        | 1.1     | 1.2      | N/A    | N/A   |                                                                                  |
| TMC Provisioned    | Yes     | No       | N/A    | N/A   | TMC-provisioned clusters implement PSPs which have an affect on app deployments. |
| HA Control Plane   | Yes     | No       | N/A    | N/A   |                                                                                  |
| CNI                | Calico  | Antrea   | Calico | N/A   |                                                                                  |
| Minimal Resources  | No      | Yes      | N/A    | N/A   | VMware clusters should test the minimal deployment size.                         |
| DNS Provider       | route53 | internal | xip.io | N/A   | Different DNS providers as per the `tanzu_dns_provider` option are tested.       |

## Stages

The following describes the stages and their purpose within the CI/CD pipeline:

- **Lint:** ensure things like YAML, file structures, directory structures all look the same.  This gives the code base a feeling like one person wrote
it and provides a level of consistency that makes all code easy to read, regardless of who wrote it.
- **Unit:** this tests the individual modules, along with their minimal dependencies, only when code within that module changes.  This lets us verify
quickly if a code has broke a change within a module.  **NOTE:** due to performance reasons, this has been temporarily disabled.
- **E2E:** this tests a full deployment of the module.  This only happens in the context of a merge request that has the `WIP` flag removed.

## Environment

We have temporarily procured a Pez environment to setup our testing infrastructure with.  All of the testing infrastructure is
run on virtual machines which host our Gitlab Runners.  We eventually would like to run this in Kubernetes, but we wanted to run
tests using Docker-in-Docker and that made for some challenges inside a Kubernetes environment.

The following environment is utilized for CI/CD testing:

| Stage   | Cluster/Host             | Runner Tags    | Executor | Runs                                                             | Comments                                                                        |
| ------- | ------------------------ | -------------- | -------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------- |
| Lint    | docker-host-generic-1    | generic        | shell    | Always                                                           | Tested with shell executor to better simulate user workflow.                    |
| Unit    | docker-host-unit-{1..5}  | unit           | shell    | Merge Requests without `WIP` in commit title                     | Tested with KIND.  Turned off for now due to performance problems with Pez Lab. |
| E2E     | docker-host-<PROVIDER>-1 | <PROVIDER>-e2e | shell    | Merge Requests with `WIP` in commit title or Merge Request Title |                                                                                 |
| Release | *                        | None           | k8s      | On tagged release of repo                                        | Runs on any available runner.                                                   |

Also see [../ci/runners/README.md](../ci/runners/README.md) for more information on runner setup.
