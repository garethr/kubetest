#!/usr/bin/env bats

@test "Fail when run with basic examples" {
  run kubetest --tests examples/tests examples/rc.yaml
  [ "$status" -eq 1 ]
}
