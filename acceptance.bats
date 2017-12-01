#!/usr/bin/env bats

@test "Fail when run with basic examples" {
  run kubetest --tests examples/tests examples/rc.yaml
  [ "$status" -eq 1 ]
  [ "$output" = "WARN examples/rc.yaml ReplicationController should have at least 4 replicas" ]
}
