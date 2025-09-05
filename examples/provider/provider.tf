# Copyright (c) HashiCorp, Inc.

# Configuration-based authentication
provider "ranger" {
  username = "admin"
  password = "rangerR0cks!"
  host     = "http://localhost:6080"
}
