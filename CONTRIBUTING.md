# Contributing to terraform-provider-ranger

We're excited to have you contribute to our project! This document outlines the process you should follow to contribute effectively.

## How to Contribute

Basic instructions about where to send patches, check out source code, and get development support:

- **Patches**: Please raise patches as Pull Requests [here](https://github.com/G-Research/terraform-provider-ranger/pulls).
- **Source Code**: You can check out the source code at https://github.com/G-Research/terraform-provider-ranger.git.
- **Support**: For development support, please reach out by [raising an issue](https://github.com/G-Research/terraform-provider-ranger/issues).

## Getting Started

Before you start contributing, please follow these steps:

- **Installation Steps**: Clone this repository.
- **Pre-requisites**: Go v1.24+, Terraform CLI.
- **Working with Source Code**: Use your preferred IDE or text editor. See [Building the Project](#building-the-project) for build instructions and [Testing Conventions](#testing-conventions) for running tests.

## Team

Understand our team structure and guidelines:

- For details on our team and roles, please see the [MAINTAINERS.md](MAINTAINERS.md) file.

## Building the Project

Ensure you can build the project successfully:

- **Build Scripts/Instructions**: Ensure this project can be built by running `go build`. Install the provider locally by following the instructions in the [Terraform documentation](https://www.terraform.io/docs/cli/config/config-file.html#third-party-plugins):

1. Create or update dev overrides in `~/.terraformrc` (or `terraform.rc` on Windows) to point at your $GOPATH:

    ```hcl
    provider_installation {
      dev_overrides {
        "registry.terraform.io/gresearch/ranger" = "/home/user/go/bin"
      }
      direct {}
    }
    ```

2. Install the provider by running:

    ```bash
    go install
    ```

3. Verify the installation by creating a `main.tf` file with the following content:

    ```hcl
    terraform {
      required_providers {
        ranger = {
          source = "registry.terraform.io/gresearch/ranger"
          version = ">= 0.1.0"
        }
      }
    }

    provider "ranger" {
      username = "admin"
      password = "rangerR0cks!"
      host     = "http://localhost:6080"
    }
    ```

4. Run `terraform plan` in the directory containing `main.tf` and verify no errors.

## Workflow and Branching

Our preferred workflow and branching structure:

- We recommend using [git flow](https://nvie.com/posts/a-successful-git-branching-model/).

## Testing Conventions

Our approach to testing:

- **Test Location**: Tests can be found in the [internal/provider](/internal/provider) directory.
- **Running Tests**: For unit tests run `make test`. For acceptance tests run `make testacc`.

## Coding Style and Linters

Our coding standards and tools:

- **Coding Standards**: Code should be formatted using `go fmt`.
- **Linters**: We use [golangci-lint](https://github.com/golangci/golangci-lint) to lint this codebase.

## Writing Issues

Guidelines for filing issues:

- **Where to File Issues**: Please file issues [here](https://github.com/G-Research/terraform-provider-ranger/issues).
- **Issue Conventions**: Please use one of the available templates when raising an issue.

## Writing Pull Requests

Guidelines for pull requests:

- **Where to File Pull Requests**: Submit your pull requests by forking this repository and raising [here](https://github.com/G-Research/terraform-provider-ranger/compare).
- **PR Conventions**: Please use the template shown when raising a Pull Request to help reviewers understand the intent behind your contribution.

## Reviewing Pull Requests

How we review pull requests:

- **Review Process**: Maintainers will be notified of new pull requests and will be reviewed periodically.
- **Reviewers**: Our reviews are conducted by [our maintainers](MAINTAINERS.md).

## Shipping Releases

Our release process:

- **Cadence**: We ship releases when new features or bug fixes are available.
- **Responsible Parties**: Releases are managed by [our maintainers](MAINTAINERS.md).

## Documentation Updates

How we handle documentation:

- **Documentation Location**: Our documentation can be found in this project's [README](README.md).
- **Update Process**: Documentation is updated whenever significant changes are made.
