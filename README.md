# Terraform Provider: Snowflake

**Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [team-cloud-foundation-tools-dl@snowflake.com](mailto:team-cloud-foundation-tools-dl@snowflake.com).

----

![.github/workflows/ci.yml](https://github.com/Snowflake-Labs/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider plugin for managing [Snowflake](https://www.snowflake.com/) accounts.

## Getting Help

If you need help, try the [discussions area](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions) of this repo.

## Install

The easiest way is to run this command:

```shell
curl https://raw.githubusercontent.com/chanzuckerberg/terraform-provider-snowflake/main/download.sh | bash -s -- -b $HOME/.terraform.d/plugins
```

**Note that this will only work with recent releases, for older releases, use the version of download.sh that corresponds to that release (replace main in that curl with the version).**

It runs a script generated by [godownloader](https://github.com/goreleaser/godownloader) which installs into the proper directory for terraform (~/.terraform.d/plugins).

You can also just download a binary from our [releases](https://github.com/Snowflake-Labs/terraform-provider-snowflake/releases) and follow the [Terraform directions for installing 3rd party plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

### Upgrading from CZI Provider
As of (5/25/2022) to provider has been transferred from CZI to Snowflake-Labs. To upgrade from CZI, please run the following command:

```shell
terraform state replace-provider chanzuckerberg/snowflake Snowflake-Labs/snowflake
```

You should also update your lock file / Terraform provider version pinning. From the deprecated source:

```hcl
# deprecated source
terraform {
  required_version = ">= 1.1.7"

  required_providers {
    snowflake = {
      source  = "chanzuckerberg/snowflake"
      version = "0.36.0"
    }
  }
}
```

To new source:

```hcl
# new source
terraform {
  required_version = ">= 1.1.7"

  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "0.36.0"
    }
  }
}
```
If you are not pinning your provider versions, you may find it useful to forcefully upgrade providers using the command: 
```
terraform init -upgrade
```

>**Note**:  0.34 is the first version published after the transfer. When the provider was transferred over not all releases were transferred for some reason. Only versions 0.28 and newer were transferred.

### For Terraform v0.13+ users

> We are now (7/29/2021) using Terraform 0.13 for testing purposes due to an issue for data sources for versions <0.13. Related PR for this change [here](https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/622#issuecomment-888879621).

You can use [Explicit Provider Source Locations](https://www.terraform.io/upgrade-guides/0-13.html#explicit-provider-source-locations).

The following maybe work well.

```terraform
terraform {
  required_providers {
    snowflake = {
      source = "Snowflake-Labs/snowflake"
      version = "0.36.0"
    }
  }
}
```

## Usage

An [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) is available from Snowflake.

In-depth docs are available [on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest).

## Development

If you do not have Go installed:

1. Install Go `brew install golang`
2. Make a Go development directory wherever you like `mkdir go_projects`
3. Add the following config to your profile

   ```shell
   export GOPATH=$HOME/../go_projects # edit with your go_projects dir
   export PATH=$PATH:$GOPATH/bin
   ```

4. Fork this repo and clone it into `go_projects`
5. cd to `terraform-provider-snowflake` and install all the required packages with `make setup`
6. Finally install goimports with `(cd && go install golang.org/x/tools/cmd/goimports@latest)`.
7. You should now be able to successfully run the tests with `make test`

It has not been tested on Windows, so if you find problems let us know.

If you want to build and test the provider locally there is a make target `make install-tf` that will build the provider binary and install it in a location that terraform can find.

## Testing

**Note: PRs for new resources will not be accepted without passing acceptance tests.**

For the Terraform resources, there are 3 levels of testing - internal, unit and acceptance tests.

The 'internal' tests are run in the `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources` package so that they can test functions that are not exported. These tests are intended to be limited to unit tests for simple functions.

The 'unit' tests are run in  `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources_test`, so they only have access to the exported methods of `resources`. These tests exercise the CRUD methods that on the terraform resources. Note that all tests here make use of database mocking and are run locally. This means the tests are fast, but are liable to be wrong in subtle ways (since the mocks are unlikely to be perfect).

You can run these first two sets of tests with `make test`.

The 'acceptance' tests run the full stack, creating, modifying and destroying resources in a live snowflake account. To run them you need a snowflake account and the proper authentication set up. These tests are slower but have higher fidelity.

To run all tests, including the acceptance tests, run `make test-acceptance`.

### Pull Request CI

Our CI jobs run the full acceptence test suite, which involves creating and destroying resources in a live snowflake account. Github Actions is configured with environment variables to authenticate to our test snowflake account. For security reasons, those variables are not available to forks of this repo.

If you are making a PR from a forked repo, you can create a new Snowflake Enterprise trial account and set up Travis to build it by setting these environment variables:

* `SNOWFLAKE_ACCOUNT` - The account name
* `SNOWFLAKE_USER` - A snowflake user for running tests.
* `SNOWFLAKE_PASSWORD` - Password for that user.
* `SNOWFLAKE_ROLE` - Needs to be ACCOUNTADMIN or similar.
* `SNOWFLAKE_REGION` - Default is us-west-2, set this if your snowflake account is in a different region.

You will also need to generate a Github API token and add the secret:

* `REVIEWDOG_GITHUB_API_TOKEN` - A token for reviewdog to use to access your github account with privileges to read/write discussion.

## Releasing

## Running a release

**Note: releases can only be done by those with keybase pgp keys allowed in the terraform registry.**

Releases will be performed once a week on **Monday around 11am PST**. If your change is more urgent and you need to use it sooner, use the commit hash.

Releases are done by [goreleaser](https://goreleaser.com/) and run by our make files. There two goreleaser configs, `.goreleaser.yml` for regular releases and `.goreleaser.prerelease.yml` for doing prereleases (for testing).

Releases are [published to the terraform registry](https://registry.terraform.io/providers/chanzuckerberg/snowflake/latest), which requires that releases by signed.

## Adding a new releaser

To set up a new person for releasing, there are a few steps–

1. releaser: a [keybase account](https://keybase.io/) and a workstation set up with their [Keybase app](https://keybase.io/download).
2. releaser: a pgp key - `keybase pgp gen`
3. releaser: export public key.
   1. If you have a single key in keybase–
      1. `keybase pgp export`
   2. If you have more than one key–
      1. `keybase pgp export` to find id if key you want to export
      2. `keybase pgp export -q KEY_ID`
4. github admin for chanzuckerberg: take public key exported above and add it [in the registry](https://registry.terraform.io/settings/gpg-keys)
5. releaser: set `KEYBASE_KEY_ID` environment variable. Note that this is different from the previous id. Get this one from `keybase pgp list`. It should be like ~70 characters long.
6. set `GITHUB_TOKEN` environment variable with a personal access token
7. releaser: run `make release-prerelease` to test that releases are working correctly
8. releaser: run `make release` to release for real
