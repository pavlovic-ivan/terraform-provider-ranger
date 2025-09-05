# terraform-provider-ranger

A Terraform provider for Apache Ranger.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

See the [examples](examples) directory for provider, resource, and data source usage.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, first start up a local Ranger instance using docker:

```
docker compose -f docker/docker-compose.yml up -d
```

Then run `make testacc` to run the acceptance tests.


```shell
make testacc
```

## Contributing

If you would like to contribute to this project, please fork the repository and make your changes in a separate branch. Once you have made your changes, you can submit a pull request for review. See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## Security

Please see [SECURITY.md](SECURITY.md) for our Security and Coordinated Vulnerability Disclosure Policy.

## Code of Conduct

This project adheres to our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

This project is licensed under the Apache License, Version 2.0 - see the [LICENSE](LICENSE) file for details.

**SPDX-License-Identifier:** Apache-2.0
