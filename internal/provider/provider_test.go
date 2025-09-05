// Copyright (c) HashiCorp, Inc.

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Ranger client is properly configured.
	// It is also possible to use the RANGER_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "ranger" {
  username = "admin"
  password = "rangerR0cks!"
  host     = "http://localhost:6080"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ranger": providerserver.NewProtocol6WithError(New("test")()),
	}
)
