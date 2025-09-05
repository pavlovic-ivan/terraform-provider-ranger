// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPolicyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ranger_policy" "test" {
  name = "test-policy"
  description = "policy description"
  service = "dev_kafka"
    resources = {
    topic = {
      values = [
        "topic-pattern-0-*",
        "topic-pattern-1-*",
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
        "example-user",
      ]
	  groups = [
		"example-group",
	  ]
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values set in the config
					resource.TestCheckResourceAttr("ranger_policy.test", "service", "dev_kafka"),
					resource.TestCheckResourceAttr("ranger_policy.test", "name", "test-policy"),
					resource.TestCheckResourceAttr("ranger_policy.test", "description", "policy description"),

					// Verify Topic resource
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.values.0", "topic-pattern-0-*"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.values.1", "topic-pattern-1-*"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.is_excludes", "false"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.is_recursive", "false"),

					// Verify policy items
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.accesses.0.type", "publish"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.accesses.0.is_allowed", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.users.0", "example-user"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.groups.0", "example-group"),

					// Verify default values
					resource.TestCheckResourceAttr("ranger_policy.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_type", "0"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_priority", "0"),
					resource.TestCheckResourceAttr("ranger_policy.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "service_type", "kafka"),
					resource.TestCheckResourceAttr("ranger_policy.test", "is_deny_all_else", "false"),

					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("ranger_policy.test", "id"),
					resource.TestCheckResourceAttrSet("ranger_policy.test", "guid"),
					resource.TestCheckResourceAttrSet("ranger_policy.test", "version"),
				),
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "ranger_policy" "test" {
  name = "test-policy"
  description = "policy description has been updated"
  service = "dev_kafka"
    resources = {
    topic = {
      values = [
        "topic-pattern-0-*",
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
        "example-user",
		"second-user",
      ]
	  groups = [
		"example-group-modified",
	  ]
    },
	{
      accesses = [
        {
          type = "consume",
        }
      ]
	  groups = [
		"consumption-group",
	  ]
    }
  ]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify values set in the config
					resource.TestCheckResourceAttr("ranger_policy.test", "service", "dev_kafka"),
					resource.TestCheckResourceAttr("ranger_policy.test", "name", "test-policy"),
					resource.TestCheckResourceAttr("ranger_policy.test", "description", "policy description has been updated"),

					// Verify Topic resource
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.values.0", "topic-pattern-0-*"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.values.1", "topic-pattern-1-*"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.values.2", "topic-pattern-2-*"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.is_excludes", "false"),
					resource.TestCheckResourceAttr("ranger_policy.test", "resources.topic.is_recursive", "false"),

					// Verify policy items
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.accesses.0.type", "publish"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.accesses.0.is_allowed", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.users.0", "example-user"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.users.1", "second-user"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.0.groups.0", "example-group-modified"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.1.accesses.0.type", "consume"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.1.accesses.0.is_allowed", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_items.1.groups.0", "consumption-group"),

					// Verify default values
					resource.TestCheckResourceAttr("ranger_policy.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_type", "0"),
					resource.TestCheckResourceAttr("ranger_policy.test", "policy_priority", "0"),
					resource.TestCheckResourceAttr("ranger_policy.test", "is_enabled", "true"),
					resource.TestCheckResourceAttr("ranger_policy.test", "service_type", "kafka"),
					resource.TestCheckResourceAttr("ranger_policy.test", "is_deny_all_else", "false"),

					// Verify dynamic values have any value set in the state
					resource.TestCheckResourceAttrSet("ranger_policy.test", "id"),
					resource.TestCheckResourceAttrSet("ranger_policy.test", "guid"),
					resource.TestCheckResourceAttrSet("ranger_policy.test", "version")),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
