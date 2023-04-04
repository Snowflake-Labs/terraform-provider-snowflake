# Terraform Provider: Snowflake

**Please note**: If you believe you have found a security issue, _please responsibly disclose_ by contacting us at [team-cloud-foundation-tools-dl@snowflake.com](mailto:team-cloud-foundation-tools-dl@snowflake.com).

----

![.github/workflows/ci.yml](https://github.com/Snowflake-Labs/terraform-provider-snowflake/workflows/.github/workflows/ci.yml/badge.svg)

This is a terraform provider plugin for managing [Snowflake](https://www.snowflake.com/) accounts.

## Getting Help

If you need help, try the [discussions area](https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions) of this repo. We also use this forum to discuss new features and changes to the provider.

**Note**: If you are an enterprise customer, please contact your Snowflake account representative. We prioritize support over GitHub issues. Also it helps us with allocating additional engineering resources to supporting the provider.

## Install

You can install the provider using `terraform init`, all you need to do is include the following block in your Terraform settings configuration. Refer to [Explicit Provider Source Locations](https://www.terraform.io/upgrade-guides/0-13.html#explicit-provider-source-locations) for more information.

```hcl
terraform {
  required_providers {
    snowflake = {
      source  = "Snowflake-Labs/snowflake"
      version = "~> 0.60"
    }
  }
}
```

### Upgrading from CZI Provider
As of (5/25/2022) to provider has been transferred from Chan Zuckerberg Initiative (CZI) to Snowflake-Labs. To upgrade from CZI, please run the following command:

```shell
terraform state replace-provider chanzuckerberg/snowflake Snowflake-Labs/snowflake
```

You should also update your lock file / Terraform provider version pinning. From the deprecated source:

```hcl
# deprecated source
terraform {
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

>**Note**:  0.34 is the first version published after the transfer. When the provider was transferred over not all of the older releases were transferred for some reason. Only versions 0.28 and newer were transferred. If you are using a version older than 0.28, it is highly recommended to upgrade to a newer version.

## Usage

An [introductory tutorial](https://guides.snowflake.com/guide/terraforming_snowflake/#0) is available from Snowflake.

In-depth docs are available [on the Terraform registry](https://registry.terraform.io/providers/Snowflake-Labs/snowflake/latest).

## Development

If you do not have Go installed:

1. Install Go `brew install golang`
2. Ensure that your GOPATH is set correctly
3. Fork this repo and clone it into `~/go/src/github.com/Snowflake-Labs/terraform-provider-snowflake`
4. cd to `terraform-provider-snowflake` and install all the required packages with `go get`
5. Build provider with `go install`

## Testing
The following environment variables need to be set for acceptance tests to run:
* `SNOWFLAKE_ACCOUNT` - The account name
* `SNOWFLAKE_USER` - A snowflake user for running tests.
* `SNOWFLAKE_PASSWORD` - Password for that user.
* `SNOWFLAKE_ROLE` - Needs to be ACCOUNTADMIN or similar.
* `SNOWFLAKE_REGION` - Default is us-west-2, set this if your snowflake account is in a different region.
* `TF_ACC` - to enable acc tests.

e.g.

```
export SNOWFLAKE_ACCOUNT=TESTACCOUNT
export SNOWFLAKE_USER=TEST_USER
export SNOWFLAKE_PASSWORD=hunter2
export SNOWFLAKE_ROLE=ACCOUNTADMIN
export SNOWFLAKE_REGION=us-west-2
export TF_ACC=true
```

**Note: PRs for new resources will not be accepted without passing acceptance tests.**

For the Terraform resources, there are 3 levels of testing - internal, unit and acceptance tests.

The 'internal' tests are run in the `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources` package so that they can test functions that are not exported. These tests are intended to be limited to unit tests for simple functions.

The 'unit' tests are run in  `github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources_test`, so they only have access to the exported methods of `resources`. These tests exercise the CRUD methods that on the terraform resources. Note that all tests here make use of database mocking and are run locally. This means the tests are fast, but are liable to be wrong in subtle ways (since the mocks are unlikely to be perfect).

You can run these first two sets of tests with `make test`.

The 'acceptance' tests run the full stack, creating, modifying and destroying resources in a live snowflake account. To run them you need a snowflake account and the proper authentication set up. These tests are slower but have higher fidelity.

To run all tests, including the acceptance tests, run `make test-acceptance`.


If you are making a PR from a forked repo, you can create a new Snowflake Enterprise trial account and set up Travis to build it by setting these environment variables:

## Advanced Debugging
If you want to build and test the provider locally there is a make target `make install-tf` that will build the provider binary and install it in a location that terraform can find.

To debug the provider with a debugger:
1. Launch the provider with the `-debug` command line argument in your debugger session. Once the provider starts, it will print instructions on setting the `TF_REATTACH_PROVIDERS` environment variable.
   ```
   Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

   Command Prompt:	set "TF_REATTACH_PROVIDERS={"registry.terraform.io/Snowflake-Labs/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr": {"Network":"tcp","String":"127.0.0.1:54706"}}}"
   PowerShell:	$env:TF_REATTACH_PROVIDERS='{"registry.terraform.io/Snowflake-Labs/snowflake":{"Protocol":"grpc","ProtocolVersion":5,"Pid":35140,"Test":true,"Addr":{"Network":"tcp","String":"127.0.0.1:54706"}}}'
   ```
2. Open a terminal where you will execute Terraform and set the `TF_REATTACH_PROVIDERS` environment variable using the command from the first step.
3. Run Terraform as usual from this terminal. Any breakpoints you set will halt execution and you can troubleshoot the provider from your debugger.

**Note**: The `TF_REATTACH_PROVIDERS` environment variable needs to be set every time you restart your debugger session as some values like the `Pid` or the TCP port will change with every execution.

For further instructions, please check the official [Terraform Plugin Development guide](https://www.terraform.io/plugin/debugging#starting-a-provider-in-debug-mode).

## Contributing

We use [Conventional Commits](https://www.conventionalcommits.org/) for commit messages formatting. Please try to adhere to the standard.
Validation is done with this regular expression:

https://github.com/Snowflake-Labs/terraform-provider-snowflake/blob/main/.github/workflows/titleLint.yml#L17

## Releasing

Releases will be performed as needed, typically once every 1-2 weeks. If your change is more urgent and you need to use it sooner, use the commit hash.

Releases are done by [goreleaser](https://goreleaser.com/) and run by our make files. There two goreleaser configs, `.goreleaser.yml` for regular releases and `.goreleaser.prerelease.yml` for doing prereleases (for testing).

Releases are [published to the terraform registry](https://registry.terraform.io/providers/chanzuckerberg/snowflake/latest), which requires that releases by signed.
