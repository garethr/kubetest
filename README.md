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


`kubetest` is currently alpha quality and undoutedly has a few issues. Please open issues if you have feedback when trying it out.


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

## Thanks

A big thank you goes to the authors of [stretchr/testify](https://github.com/stretchr/testify/) from where much of the assertion code has been ported. 

