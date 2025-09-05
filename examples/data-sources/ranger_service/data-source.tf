# Copyright (c) HashiCorp, Inc.

# Get the dev kafka policy.
data "ranger_policy" "dev_kafka" {
  name = "dev_kafka"
}
