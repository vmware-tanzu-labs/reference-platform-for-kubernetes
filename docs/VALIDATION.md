# Validation Tests

The objective of RPK validation tests is to provide automated
testing of modules' end-user functionality.  These are implemented as sonobuoy
plugins.

## Running Tests

For a module that has validation tests, run them with:

    $ ROLE=workload-tenancy make validate.role

To view the status of the validation test:

    $ make validate.status

Once status is `complete`, if the test result is `failed`, see Test Failures
section below.  If test result is `passed`, clean up the sonobuoy assets as
follows:

    $ make validate.clean

## Test Failures

In the event a test fails you will need to retrieve the test results to gather
more information:

    $ make validate.retrieve

This will download and store the results locally in
`/tmp/rpk-validation-results`.  To view a summary of results:

    $ make validate.results

The results will display how many tests passed and failed.  In order to get more
detail about what each test validates:

    $ make validate.detail

If you need log output from the plugin:

    $ make validate.logs

## Test Development

Tests can be written in any language.  Prefer languages that are commonly used
by others on the team and have good client libraries for the operations being
tested.  If in doubt, stick with Go or Python.  Bash can be acceptable for
simple test cases, too.

If developing tests for a module the process roughly is:

1. Determine _what_ to test.  Think about the manual steps that you would go
   through to validate that a module is functioning properly.
2. Write a program to conduct those steps.  If the steps involve creating and
   inspecting Kubernetes resources, use a Kubernetes client library to interact
   with the Kubernetes API.  If the steps involve checking services in with your
   web browser, use HTTP libraries to call the same endpoints and check the
   responses.
3. Containerized the program so that it may be run in the cluster being
   validated.
4. Create a [sonobuoy](https://sonobuoy.io/) plugin configuration to run the
   tests.

Plugin code lives in each modules' `validation` directory.

Check out existing validation tests as well as the sonobuoy repo for [example
plugins](https://github.com/vmware-tanzu/sonobuoy/tree/v0.19.0/examples/plugins).

## Conventions

Sonobuoy plugin config filename: `sonobuoy-plugin.yaml`.  This allows users to
run the validation test for a module using make targets.

By default, log output from your validation tests to a file:
`/tmp/results/plugin.log`.

