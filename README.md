# Kubetest

`kubetest` is a tool for running tests against a Kubernetes YAML or JSON configuration file.
These tests can be used to enforce local or global best-practices, for example:

* Ensuring certain labels are set
* Prevent usage of images with the `latest` tag
* Prohibit privileged containers
* Enforce a naming convention for different resources

[![Build
Status](https://travis-ci.org/garethr/kubetest.svg)](https://travis-ci.org/garethr/kubetest)
[![Go Report
Card](https://goreportcard.com/badge/github.com/garethr/kubetest)](https://goreportcard.com/report/github.com/garethr/kubetest)
[![GoDoc](https://godoc.org/github.com/garethr/kubetest?status.svg)](https://godoc.org/github.com/garethr/kubetest)
[![Coverage
Status](https://coveralls.io/repos/github/garethr/kubetest/badge.svg?branch=master)](https://coveralls.io/github/garethr/kubetest?branch=master)


`kubetest` is currently alpha quality and undoutedly has a few issues. Things will change, hopefully for the better. Please open issues if you have feedback when trying it out.


## Writing tests

Tests are written in [Skylark](https://github.com/google/skylark), which is a small dialect of Python suitable for embedding in other programmes. This means you do not need an additional interpreter installed to run tests with `kubetest`. `kubetest` prioritises interopability over flexibility in this regard. Tests for Kubetest just require the `kubetest` binary to run. Let's take a look at an example test:

```python
#// vim: set ft=python:
def test_for_team_label():
    if spec["kind"] == "Deployment":
        labels = spec["spec"]["template"]["metadata"]["labels"]:
        assert_contains(labels, "team", "should indicate which team owns the deployment")

test_for_team_label()
```

Save the test file in a directory called `tests`, with an extension of `.sky`. You can change the default directory name using the `--tests` flag. You can now run `kubetest` against your configuration files.

```bash
$ kubetest my-deployment.yaml
WARN my-deployment.yaml Deployment should have at least 4 replicas
$ echo $?
1
``` 

If any of the tests fail then `kubetest` will return a non-zero exit code.

By default `kubetest` outputs information about failing tests only, but you can pass `--verbose` to get information about passing tests as well.

```bash
$ kubetest rc.yaml --verbose 
INFO rc.yaml should not use latest images
WARN rc.yaml ReplicationController should have at least 4 replicas
```

## The `spec` variable

`spec` is a global variable passed into the Skylark code which contains the structure of the Kubernetes configuration passed in to `kubetest`. You'll need to be reasonably familiar with the structure of the Kubernetes API objects to write tests, but it is possible to write helper methods for common assertions.


## Assertions

`kubetest` automatically makes available a set of assertions to make writing tests in Skylark more pleasant. A failed assertion results in `kubetest` exiting with a non-zero exit code, and assertions output results as shown above.

* assert_equal
* assert_contains
* assert_not_contains
* assert_not_equal
* assert_nil
* assert_not_nil
* fail
* fail_now
* assert_empty
* assert_not_empty
* assert_true
* assert_false

Assertions take zero, one or two arguments (noted above) depending on what they are comparing. They then take an additional message argument which is output when the assertion runs. For example the following assertion checks whether the variable `labels` contains the value `team`.

```python
assert_contains(labels, "team", "should indicate which team owns the deployment")
```


## Installation

Tagged versions of `kubetest` are built by Travis and automatically
uploaded to GitHub. This means you should find `tar.gz` files under the
release tab. These should contain a single `kubetest` binary for platform
in the filename (ie. windows, linux, darwin). Either execute that binary
directly or place it on your path.

```
wget https://github.com/garethr/kubetest/releases/download/0.1.0/kubetest-darwin-amd64.tar.gz
tar xf kubetest-darwin-amd64.tar.gz
cp kubetest /usr/local/bin
```

Windows users can download tar or zip files from the releases tab.


## CLI

```
$ kubetest --help
Run tests against a Kubernetes YAML file

Usage:
  kubetest <file> [file...] [flags]

Flags:
  -h, --help           help for kubetest
      --json           Output results as JSON
  -t, --tests string   Test directory (default "tests")
      --verbose        Output passes as well as failures
```


## Thanks

A big thank you goes to the authors of [stretchr/testify](https://github.com/stretchr/testify/) from where much of the assertion code has been ported. 

