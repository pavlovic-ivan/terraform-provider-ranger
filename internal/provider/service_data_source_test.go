// Copyright (c) HashiCorp, Inc.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + `data "ranger_service" "test" { name = "dev_kafka" }`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify id and name are set
					resource.TestCheckResourceAttr("data.ranger_service.test", "id", "6"),
					resource.TestCheckResourceAttr("data.ranger_service.test", "name", "dev_kafka"),
				),
			},
		},
	})
}
