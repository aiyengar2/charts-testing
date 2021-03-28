## Unit Testing for Helm Chart

### What is this directory?

This directory contains Go tests that utilize [rancher/kube-tester](https://github.com/rancher/kube-tester), a unit testing framework for Kubernetes manifests based on [stackrox/kube-linter](https://github.com/stackrox/kube-linter), to run unit tests on Helm charts.

### How do I run a test?

```bash
$ go test -v packages/<package>/tests -strict=<true|false> -chart=<path-to-chart>
```

If the flags are omitted, `-strict=true` and `-chart=packages/<package>/tests` are the default values.

### How do I add a test?

Adding a test involves two actions:

1. Writing a test function that operates on sets of k8s objects

2. Calling the test function in `packages/<package>/tests/<test-group-name>_test.go` on specific templates

### Writing a test

TBD...

### Adding a test to a test suite

TBD...

### How do I add a new template to run tests on?

Drop the values.yaml file in `packages/<package>/tests/values/<template-name>.yaml`.

### Important Files and Directories

- `values/`: A directory that will contain a set of values.yaml for your charts. These files are automatically loaded as templates on the testing suite on running any test (note: the file `values/default.yaml` corresponds to the template `default.yaml`).

- `main_test.go`: A file that sets up the testing framework. **This file should never be modified.**

- `package_test.go`: A file that contains package-specific logic for the testing framework. **This file can be modified.**

### Working with Custom Resource Definitions

While working with this framework, you might encounter an error like:

```bash
* failed to decode: no kind "ClusterScanBenchmark" is registered for version "cis.cattle.io/v1" in scheme "pkg/runtime/scheme.go:100"
```

This error shows up because the scheme used by the testing framework does not recognize the custom resource definition that appears in your template. To resolve this problem, modify `package_test.go` to call the relevant `AddToScheme` function, which should be exposed by the operator

### Maintainers
- aiyengar2 (arvind.iyengar@suse.com)
- rancher-max (max.ross@suse.com)
- brendarearden (brenda.rearden@suse.com)