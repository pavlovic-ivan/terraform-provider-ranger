# Copyright (c) HashiCorp, Inc.

# Kafka policy
resource "ranger_policy" "kafka_policy" {
  service = "dev_kafka"
  name    = "Example Kafka Policy"
  resources = {
    topic = {
      values = [
        "topic-pattern-1-*",
        "topic-pattern-2-*",
      ]
    }
  }
  policy_items = [
    {
      accesses = [
        {
          type = "publish",
        }
      ]
      users = [
        "user1",
        "user2",
      ]
      groups = [
        "group1",
      ]
    },
    {
      accesses = [
        {
          type = "consume",
        }
      ]
      users = [
        "user1",
        "user2",
      ]
      groups = [
        "group1",
      ]
    },
    {
      accesses = [
        {
          type = "create",
        }
      ]
      users = [
        "admin1",
      ]
    },
    {
      accesses = [
        {
          type = "delete",
        }
      ]
      users = [
        "admin1",
      ]
    }
  ]
}
